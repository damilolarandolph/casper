package main

import "github.com/damilolarandolph/casper/cpu"

func main() {
	cpu := cpu.NewArm7(2)
	cpu.Tick()
	cpu.Run()
}

func init() {
}
