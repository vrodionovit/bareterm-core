package main

import (
	"fmt"
)

func (t *Terminal) moveCursorUp(params string) {
	fmt.Printf("moveCursorUp called with params: %s\n", params)
}

func (t *Terminal) moveCursorDown(params string) {
	fmt.Printf("moveCursorDown called with params: %s\n", params)
}

func (t *Terminal) moveCursorForward(params string) {
	fmt.Printf("moveCursorForward called with params: %s\n", params)
}

func (t *Terminal) moveCursorBackward(params string) {
	fmt.Printf("moveCursorBackward called with params: %s\n", params)
}

func (t *Terminal) setGraphicsMode(params string) {
	fmt.Printf("setGraphicsMode called with params: %s\n", params)
	t.handleColor(params)
}

func (t *Terminal) eraseLine(params string) {
	fmt.Printf("eraseLine called with params: %s\n", params)
}

func (t *Terminal) scrollUp(params string) {
	fmt.Printf("scrollUp called with params: %s\n", params)
}

func (t *Terminal) scrollDown(params string) {
	fmt.Printf("scrollDown called with params: %s\n", params)
}

func (t *Terminal) deviceStatusReport(params string) {
	fmt.Printf("deviceStatusReport called with params: %s\n", params)
}

func (t *Terminal) setMode(sequence string, enable bool) {
	fmt.Printf("setMode called with sequence: %s, enable: %v\n", sequence, enable)
}

func (t *Terminal) setScrollingRegion(params string) {
	fmt.Printf("setScrollingRegion called with params: %s\n", params)
}

func (t *Terminal) clearScreen() {
	fmt.Println("Screen cleared")
}
func (t *Terminal) moveCursor(params string) {
	fmt.Printf("moveCursor called with params: %s\n", params)
}
