package main

import "fmt"

type replica struct {
	id    int
	value string
	ver   int
}

const N, W, R = 5, 3, 3

func write(replicas []*replica, val string, ver int) int {
	acks := 0
	for _, r := range replicas {
		r.value, r.ver = val, ver
		acks++
		if acks >= W {
			break
		}
	}
	return acks
}

func read(replicas []*replica) (string, int) {
	var best string
	bestVer := -1
	collected := 0
	for _, r := range replicas {
		if r.ver > bestVer {
			bestVer, best = r.ver, r.value
		}
		collected++
		if collected >= R {
			break
		}
	}
	return best, bestVer
}

func main() {
	replicas := make([]*replica, N)
	for i := range replicas {
		replicas[i] = &replica{id: i}
	}

	// Write quorum W=3 out of 5
	acks := write(replicas, "alice", 1)
	fmt.Printf("write acks=%d (quorum=%d)\n", acks, W)

	// Leave replica 4 stale
	replicas[4].value = "stale"
	replicas[4].ver = 0

	// Read quorum R=3: read first 3 (all fresh)
	v, ver := read(replicas)
	fmt.Printf("read value=%q ver=%d\n", v, ver)

	// Read including stale replica (replicas[3] is ok, [4] is stale)
	v2, ver2 := read(replicas[2:]) // reads [2,3,4] — ver 1,1,0 — still picks ver 1
	fmt.Printf("read with stale=%q ver=%d (stale ignored by max-ver)\n", v2, ver2)
}
