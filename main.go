package main

import (
	"os"
	"flag"

	"Ataxx/ataxx"
	"Ataxx/wgui"
)


func parseInput() (p wgui.Params) {
	path := flag.String("path", "nopath", "Path to the directory. If the icons don't load the path hasn't been properly input. Eg.: /home/.../AtaxxGUI-master")
	tc := flag.String("tc", "inf", "Time control to use, default is infinite. 12+5 is 12secs and increment of 5secs")

	flag.Parse()

	p.Path = *path
	p.Tc = wgui.ParseTC(*tc)

	return
}

func init() {
	ataxx.InitAtaxx()
}

func main() {
	params := parseInput()
	wgui.Start(len(os.Args), os.Args, params)
}
