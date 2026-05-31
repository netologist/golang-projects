package main

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type ring struct {
	points []uint32
	owner  map[uint32]string
}

func newRing(nodes ...string) *ring {
	r := &ring{owner: map[uint32]string{}}
	for _, n := range nodes {
		h := crc32.ChecksumIEEE([]byte(n))
		r.points = append(r.points, h)
		r.owner[h] = n
	}
	sort.Slice(r.points, func(i, j int) bool { return r.points[i] < r.points[j] })
	return r
}

func (r *ring) get(key string) string {
	h := crc32.ChecksumIEEE([]byte(key))
	idx := sort.Search(len(r.points), func(i int) bool { return r.points[i] >= h })
	if idx == len(r.points) {
		idx = 0
	}
	return r.owner[r.points[idx]]
}

func main() {
	r := newRing("node-a", "node-b", "node-c")
	for _, key := range []string{"user:1", "user:2", "user:3", "user:4"} {
		fmt.Printf("%s -> %s\n", key, r.get(key))
	}
}
