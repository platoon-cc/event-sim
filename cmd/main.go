package main

import sim "github.com/platoon-cc/event-sim"

func main() {
	err := sim.SimulateForProject(1)
	if err != nil {
		panic(err)
	}
}
