package main

import "testing"

func TestGossip_Convergence(t *testing.T) {
	for trial := 0; trial < 10; trial++ {
		c := NewCluster(10, 2)
		c.Seed(0)
		rounds := c.ConvergedIn(200)
		if rounds == -1 {
			t.Fatal("did not converge within 200 rounds")
		}
	}
}

func TestGossip_AllInfectedAfterConvergence(t *testing.T) {
	c := NewCluster(20, 3)
	c.Seed(0)
	rounds := c.ConvergedIn(500)
	if rounds == -1 {
		t.Fatal("did not converge")
	}
	if c.InfectedCount() != c.TotalNodes() {
		t.Fatalf("not all infected: %d/%d", c.InfectedCount(), c.TotalNodes())
	}
}

func TestGossip_FanoutEffect(t *testing.T) {
	// Higher fanout should converge faster on average
	const trials = 20
	var sumLow, sumHigh int
	for i := 0; i < trials; i++ {
		cLow := NewCluster(20, 1)
		cLow.Seed(0)
		sumLow += cLow.ConvergedIn(1000)

		cHigh := NewCluster(20, 4)
		cHigh.Seed(0)
		sumHigh += cHigh.ConvergedIn(1000)
	}
	if sumLow <= sumHigh {
		t.Fatalf("expected lower fanout to take more rounds: low=%d high=%d", sumLow, sumHigh)
	}
}
