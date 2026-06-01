package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// NodeState represents a Raft node's role.
type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

func (s NodeState) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	default:
		return "Leader"
	}
}

// LogEntry is a single Raft log entry.
type LogEntry struct {
	Term    int
	Command string
}

// RequestVoteArgs is the RV RPC argument.
type RequestVoteArgs struct {
	Term         int
	CandidateID  int
	LastLogIndex int
	LastLogTerm  int
}

// RequestVoteReply is the RV RPC reply.
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// AppendEntriesArgs is the AE RPC argument.
type AppendEntriesArgs struct {
	Term         int
	LeaderID     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

// AppendEntriesReply is the AE RPC reply.
type AppendEntriesReply struct {
	Term    int
	Success bool
}

// RaftNode is a single node in the Raft cluster.
type RaftNode struct {
	mu          sync.Mutex
	ID          int
	state       NodeState
	currentTerm int
	votedFor    int
	log         []LogEntry
	commitIndex int
	peers       []*RaftNode
	deadCh      chan struct{}
	heartbeatCh chan struct{}

	// leader-only
	nextIndex  []int
	matchIndex []int
}

// NewNode creates a Raft node.
func NewNode(id int) *RaftNode {
	return &RaftNode{
		ID:          id,
		state:       Follower,
		votedFor:    -1,
		deadCh:      make(chan struct{}),
		heartbeatCh: make(chan struct{}, 1),
	}
}

// SetPeers wires the node to its peers.
func (n *RaftNode) SetPeers(peers []*RaftNode) {
	n.peers = peers
	n.nextIndex = make([]int, len(peers))
	n.matchIndex = make([]int, len(peers))
}

// Start launches the node's main loop.
func (n *RaftNode) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go n.run(wg)
}

func (n *RaftNode) isDead() bool {
	select {
	case <-n.deadCh:
		return true
	default:
		return false
	}
}

// Stop kills the node.
func (n *RaftNode) Stop() {
	select {
	case <-n.deadCh:
	default:
		close(n.deadCh)
	}
}

func (n *RaftNode) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if n.isDead() {
			return
		}
		timeout := time.Duration(200+rand.IntN(200)) * time.Millisecond
		select {
		case <-n.deadCh:
			return
		case <-n.heartbeatCh:
			// reset timer
		case <-time.After(timeout):
			n.mu.Lock()
			state := n.state
			n.mu.Unlock()
			if state != Leader {
				n.startElection()
			}
		}
	}
}

func (n *RaftNode) startElection() {
	n.mu.Lock()
	n.state = Candidate
	n.currentTerm++
	term := n.currentTerm
	n.votedFor = n.ID
	lastIdx := len(n.log) - 1
	lastTerm := 0
	if lastIdx >= 0 {
		lastTerm = n.log[lastIdx].Term
	}
	n.mu.Unlock()

	votes := 1
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, peer := range n.peers {
		wg.Add(1)
		go func(p *RaftNode) {
			defer wg.Done()
			reply := p.RequestVote(RequestVoteArgs{
				Term: term, CandidateID: n.ID,
				LastLogIndex: lastIdx, LastLogTerm: lastTerm,
			})
			if reply.VoteGranted {
				mu.Lock()
				votes++
				mu.Unlock()
			}
			if reply.Term > term {
				n.mu.Lock()
				n.becomeFollower(reply.Term)
				n.mu.Unlock()
			}
		}(peer)
	}
	wg.Wait()

	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state != Candidate || n.currentTerm != term {
		return
	}
	if votes > (len(n.peers)+1)/2 {
		n.becomeLeader()
	}
}

func (n *RaftNode) becomeFollower(term int) {
	n.state = Follower
	n.currentTerm = term
	n.votedFor = -1
}

func (n *RaftNode) becomeLeader() {
	n.state = Leader
	fmt.Printf("[term %d] node %d → LEADER\n", n.currentTerm, n.ID)
	go n.leaderLoop()
}

func (n *RaftNode) leaderLoop() {
	for !n.isDead() {
		n.mu.Lock()
		if n.state != Leader {
			n.mu.Unlock()
			return
		}
		term := n.currentTerm
		n.mu.Unlock()

		for _, peer := range n.peers {
			go n.sendAE(peer, term)
		}
		time.Sleep(60 * time.Millisecond)
	}
}

func (n *RaftNode) sendAE(peer *RaftNode, term int) {
	n.mu.Lock()
	entries := append([]LogEntry{}, n.log...)
	n.mu.Unlock()

	reply := peer.AppendEntries(AppendEntriesArgs{
		Term: term, LeaderID: n.ID,
		Entries:      entries,
		LeaderCommit: n.commitIndex,
	})
	if reply.Term > term {
		n.mu.Lock()
		n.becomeFollower(reply.Term)
		n.mu.Unlock()
	}
}

// RequestVote handles a RequestVote RPC.
func (n *RaftNode) RequestVote(args RequestVoteArgs) RequestVoteReply {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.isDead() {
		return RequestVoteReply{Term: n.currentTerm}
	}
	if args.Term > n.currentTerm {
		n.becomeFollower(args.Term)
	}
	grant := false
	if args.Term >= n.currentTerm && (n.votedFor == -1 || n.votedFor == args.CandidateID) {
		n.votedFor = args.CandidateID
		grant = true
		select {
		case n.heartbeatCh <- struct{}{}:
		default:
		}
	}
	return RequestVoteReply{Term: n.currentTerm, VoteGranted: grant}
}

// AppendEntries handles an AppendEntries RPC (heartbeat + log replication).
func (n *RaftNode) AppendEntries(args AppendEntriesArgs) AppendEntriesReply {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.isDead() {
		return AppendEntriesReply{Term: n.currentTerm}
	}
	if args.Term < n.currentTerm {
		return AppendEntriesReply{Term: n.currentTerm}
	}
	n.becomeFollower(args.Term)
	select {
	case n.heartbeatCh <- struct{}{}:
	default:
	}
	if len(args.Entries) > 0 {
		n.log = append([]LogEntry{}, args.Entries...)
	}
	return AppendEntriesReply{Term: n.currentTerm, Success: true}
}

// Submit appends a command to the leader's log and replicates it.
func (n *RaftNode) Submit(command string) bool {
	n.mu.Lock()
	if n.state != Leader {
		n.mu.Unlock()
		return false
	}
	n.log = append(n.log, LogEntry{Term: n.currentTerm, Command: command})
	term := n.currentTerm
	n.mu.Unlock()

	for _, peer := range n.peers {
		go n.sendAE(peer, term)
	}
	return true
}

// State returns the current NodeState of the node.
func (n *RaftNode) State() NodeState {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.state
}

// Log returns a copy of the node's log.
func (n *RaftNode) Log() []LogEntry {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]LogEntry{}, n.log...)
}

// Term returns the node's current term.
func (n *RaftNode) Term() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.currentTerm
}
