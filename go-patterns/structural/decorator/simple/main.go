package main

import (
	"fmt"
	"os"
)

func main() {
	// Compose: TimingWriter(CountingWriter(os.Stdout))
	counting := NewCountingWriter(os.Stdout)
	timing := NewTimingWriter(counting, "demo")

	fmt.Fprint(timing, "hello, decorator\n")
	fmt.Fprint(timing, "pattern\n")

	fmt.Printf("total bytes written to stdout: %d\n", counting.BytesWritten())
}
