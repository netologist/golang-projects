package main

import (
	"context"
	"fmt"
)

type resource struct {
	name string
	vote Vote
}

func (r *resource) Prepare(_ context.Context, _ string) (Vote, error) {
	fmt.Printf("%s votes %v\n", r.name, r.vote == VoteCommit)
	return r.vote, nil
}

func (r *resource) Commit(_ context.Context, _ string) error {
	fmt.Printf("%s committed\n", r.name)
	return nil
}

func (r *resource) Abort(_ context.Context, _ string) error {
	fmt.Printf("%s aborted\n", r.name)
	return nil
}

func main() {
	c := New(
		&resource{"db", VoteCommit},
		&resource{"cache", VoteCommit},
		&resource{"queue", VoteCommit},
	)
	fmt.Println("transaction result:", c.Execute(context.Background(), "tx-100"))
}
