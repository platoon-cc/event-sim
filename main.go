package main

import "github.com/platoon-cc/event-sim/sim"

func main() {
	err := sim.SimulateForProject(1)
	if err != nil {
		panic(err)
	}
}
