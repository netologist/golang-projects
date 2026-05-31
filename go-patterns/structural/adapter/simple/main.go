package main

func main() {
	var mp ModernPrinter = &PrinterAdapter{inner: &oldPrinter{}}
	mp.Print("hello adapter")
}
