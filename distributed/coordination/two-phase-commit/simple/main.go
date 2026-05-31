package main

import "fmt"

type participant struct {
	name    string
	canVote bool
}

func (p participant) prepare() bool { return p.canVote }

func twoPhaseCommit(parts []participant) string {
	for _, p := range parts {
		if !p.prepare() {
			return "ABORT (vote from " + p.name + ")"
		}
	}
	return "COMMIT"
}

func main() {
	all := []participant{{"db", true}, {"cache", true}, {"queue", true}}
	fmt.Println("result:", twoPhaseCommit(all))

	withAbort := []participant{{"db", true}, {"cache", false}}
	fmt.Println("result:", twoPhaseCommit(withAbort))
}
