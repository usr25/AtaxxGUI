package main

import (
	"os"

	"Ataxx/ataxx"
	"Ataxx/wgui"
)

func init() {
	ataxx.InitAtaxx()
}

//Use goroutines for multithreading when playing with engines
func main() {
	wgui.Start(len(os.Args), os.Args)
}
