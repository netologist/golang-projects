package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fns := map[string]func(context.Context) (int, error){
		"shard-1": func(_ context.Context) (int, error) { time.Sleep(10 * time.Millisecond); return 10, nil },
		"shard-2": func(_ context.Context) (int, error) { time.Sleep(15 * time.Millisecond); return 20, nil },
		"shard-3": func(_ context.Context) (int, error) { time.Sleep(5 * time.Millisecond); return 30, nil },
	}
	results := Gather(context.Background(), fns)

	total := 0
	for _, r := range results {
		fmt.Printf("%s: value=%d latency=%s err=%v\n", r.Source, r.Value, r.Latency.Round(time.Millisecond), r.Err)
		total += r.Value
	}
	fmt.Println("total:", total)
}
