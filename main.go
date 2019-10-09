package main

import (
	"os"

	"Ataxx/ataxx"
	"Ataxx/wgui"
)

func init() {
	ataxx.InitAtaxx()
}

func main() {
	wgui.Start(len(os.Args), os.Args)
}
