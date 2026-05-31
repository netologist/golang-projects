package main

import "fmt"

// Image is the subject interface.
type Image interface {
	Display()
}

// RealImage is the expensive concrete subject.
type RealImage struct{ filename string }

func (r *RealImage) Display() { fmt.Printf("Displaying %s\n", r.filename) }

func loadImage(filename string) *RealImage {
	fmt.Printf("Loading %s from disk...\n", filename)
	return &RealImage{filename: filename}
}

// ProxyImage lazily loads the RealImage on first Display.
type ProxyImage struct {
	filename string
	real     *RealImage
}

func (p *ProxyImage) Display() {
	if p.real == nil {
		p.real = loadImage(p.filename)
	}
	p.real.Display()
}
