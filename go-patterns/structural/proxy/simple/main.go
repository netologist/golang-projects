package main

func main() {
	var img Image = &ProxyImage{filename: "photo.png"}
	// Not loaded yet.
	img.Display() // loads, then displays
	img.Display() // already loaded, displays only
}
