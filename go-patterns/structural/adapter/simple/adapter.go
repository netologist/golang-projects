package main

import "fmt"

// LegacyPrinter uses an old interface.
type LegacyPrinter interface {
	PrintOld(s string) string
}

type oldPrinter struct{}

func (o *oldPrinter) PrintOld(s string) string { return "old: " + s }

// ModernPrinter is the interface our system expects.
type ModernPrinter interface {
	Print(s string)
}

// PrinterAdapter wraps a LegacyPrinter to satisfy ModernPrinter.
type PrinterAdapter struct {
	inner LegacyPrinter
}

func (a *PrinterAdapter) Print(s string) { fmt.Println(a.inner.PrintOld(s)) }
