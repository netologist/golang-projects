package main

import "fmt"

func main() {
	jobs := make(chan Job, 5)
	for i := 1; i <= 5; i++ {
		jobs <- Job{ID: i, Value: i}
	}
	close(jobs)

	total := 0
	for r := range Run(3, jobs) {
		total += r.Output
	}
	fmt.Println("sum of squares:", total)
}
