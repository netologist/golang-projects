package main

func main() {
	l1 := GetLogger()
	l2 := GetLogger()

	l1.Log("hello from l1")
	l2.Log("hello from l2")

	if l1 == l2 {
		l1.Log("same instance confirmed")
	}
}
