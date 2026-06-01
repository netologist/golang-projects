package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

const numNodes = 5

type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

type Node struct {
	mu         sync.Mutex
	id         int
	state      NodeState
	term       int
	votedFor   int
	log        []string
	peers      []*Node
	leaderID   int
	deadCh     chan struct{}
	resetTimer chan struct{}
}

func newNode(id int) *Node {
	return &Node{
		id:         id,
		state:      Follower,
		votedFor:   -1,
		leaderID:   -1,
		deadCh:     make(chan struct{}),
		resetTimer: make(chan struct{}, 1),
	}
}

func (n *Node) isDead() bool {
	select {
	case <-n.deadCh:
		return true
	default:
		return false
	}
}

func (n *Node) requestVote(candidateTerm, candidateID int) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.isDead() {
		return false
	}
	if candidateTerm > n.term {
		n.term = candidateTerm
		n.state = Follower
		n.votedFor = -1
	}
	if candidateTerm >= n.term && (n.votedFor == -1 || n.votedFor == candidateID) {
		n.votedFor = candidateID
		n.resetTimer <- struct{}{}
		return true
	}
	return false
}

func (n *Node) appendEntries(leaderTerm, leaderID int, entry string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.isDead() {
		return false
	}
	if leaderTerm < n.term {
		return false
	}
	n.term = leaderTerm
	n.state = Follower
	n.leaderID = leaderID
	if entry != "" {
		n.log = append(n.log, entry)
	}
	select {
	case n.resetTimer <- struct{}{}:
	default:
	}
	return true
}

func (n *Node) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if n.isDead() {
			return
		}
		timeout := time.Duration(150+rand.IntN(150)) * time.Millisecond
		select {
		case <-n.deadCh:
			return
		case <-n.resetTimer:
			// heartbeat or vote received, stay follower
		case <-time.After(timeout):
			n.startElection()
		}
	}
}

func (n *Node) startElection() {
	n.mu.Lock()
	n.state = Candidate
	n.term++
	currentTerm := n.term
	n.votedFor = n.id
	n.mu.Unlock()

	votes := 1
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, peer := range n.peers {
		wg.Add(1)
		go func(p *Node) {
			defer wg.Done()
			if p.requestVote(currentTerm, n.id) {
				mu.Lock()
				votes++
				mu.Unlock()
			}
		}(peer)
	}
	wg.Wait()

	n.mu.Lock()
	defer n.mu.Unlock()
	if n.state != Candidate || n.term != currentTerm {
		return
	}
	if votes > numNodes/2 {
		n.state = Leader
		n.leaderID = n.id
		fmt.Printf("[term %d] node %d elected LEADER (votes=%d)\n", n.term, n.id, votes)
		go n.sendHeartbeats()
	}
}

func (n *Node) sendHeartbeats() {
	for !n.isDead() {
		n.mu.Lock()
		if n.state != Leader {
			n.mu.Unlock()
			return
		}
		term := n.term
		n.mu.Unlock()

		for _, peer := range n.peers {
			go peer.appendEntries(term, n.id, "")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (n *Node) AppendToLog(entry string) bool {
	n.mu.Lock()
	if n.state != Leader {
		n.mu.Unlock()
		return false
	}
	n.log = append(n.log, entry)
	term := n.term
	n.mu.Unlock()

	for _, peer := range n.peers {
		peer.appendEntries(term, n.id, entry)
	}
	return true
}

func findLeader(nodes []*Node) *Node {
	for _, n := range nodes {
		if n.isDead() {
			continue
		}
		n.mu.Lock()
		is := n.state == Leader
		n.mu.Unlock()
		if is {
			return n
		}
	}
	return nil
}

func waitLeader(nodes []*Node) *Node {
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		if l := findLeader(nodes); l != nil {
			return l
		}
	}
	return nil
}

func main() {
	nodes := make([]*Node, numNodes)
	for i := range nodes {
		nodes[i] = newNode(i)
	}
	for i, n := range nodes {
		for j, p := range nodes {
			if i != j {
				n.peers = append(n.peers, p)
			}
		}
	}

	var wg sync.WaitGroup
	for _, n := range nodes {
		wg.Add(1)
		go n.run(&wg)
	}

	// Wait for initial leader
	leader := waitLeader(nodes)
	if leader == nil {
		fmt.Println("ERROR: no leader elected")
		return
	}
	fmt.Printf("initial leader: node %d\n", leader.id)

	// Append some entries
	for i := 1; i <= 3; i++ {
		leader.AppendToLog(fmt.Sprintf("entry-%d", i))
	}
	time.Sleep(100 * time.Millisecond)

	// Crash the leader
	fmt.Printf("crashing leader node %d\n", leader.id)
	close(leader.deadCh)

	// Wait for re-election
	newLeader := waitLeader(nodes)
	if newLeader == nil {
		fmt.Println("ERROR: no new leader")
		return
	}
	fmt.Printf("new leader: node %d\n", newLeader.id)

	// Append more entries
	for i := 4; i <= 5; i++ {
		newLeader.AppendToLog(fmt.Sprintf("entry-%d", i))
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\nFinal logs:")
	for _, n := range nodes {
		if !n.isDead() {
			n.mu.Lock()
			fmt.Printf("  node %d: %v\n", n.id, n.log)
			n.mu.Unlock()
		}
	}

	for _, n := range nodes {
		if !n.isDead() {
			close(n.deadCh)
		}
	}
	wg.Wait()
}
