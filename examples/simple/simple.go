package main

import (
	"fmt"
	//"math"
	"strconv"

	//rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	rvo "../../src/rvosimulator"
)

func setupScenario(sim *rvo.RVOSimulator) {
	//sim.SetTimeStep(0.25)
	//sim.SetAgentDefaults(15.0, 10, 10.0, 10.0, 1.5, 2.0, &rvo.Vector2{})
	//fmt.Printf("simu3 %v, %v\n", sim.DefaultAgent.MaxNeighbors, sim.DefaultAgent.NeighborDist)

	for i := 0; i < 11; i++{
		id, _ := sim.AddAgentPosition(&rvo.Vector2{X: 2, Y: 0})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 2, Y: 1})
		id2, _ := sim.AddAgentPosition(&rvo.Vector2{X: 1, Y: 1})
		sim.SetAgentPrefVelocity(id2, &rvo.Vector2{X: 1, Y: 3})
	}
	//a0, _ := sim.AddAgentPosition(&rvo.Vector2{X: 0, Y: 0})
	//a1, _ := sim.AddAgentPosition(&rvo.Vector2{X: 1, Y: 0})
	//a2, _ := sim.AddAgentPosition(&rvo.Vector2{X: 1, Y: 1})
	//a3, _ := sim.AddAgentPosition(&rvo.Vector2{X: 0, Y: 1})

	/*obstaclesPosition := []*rvo.Vector2{
		&rvo.Vector2{X: 0.1, Y: 0.1},
		&rvo.Vector2{X: -0.1, Y: 0.1},
		&rvo.Vector2{X: -0.1, Y: -0.1},
	}
	obstaclesPosition2 := []*rvo.Vector2{
		&rvo.Vector2{X: 0.2, Y: 0.2},
		&rvo.Vector2{X: -0.2, Y: 0.2},
	}*/
	//sim.AddObstacle(obstaclesPosition)
	//sim.AddObstacle(obstaclesPosition2)
	//sim.ProcessObstacles()

	//sim.SetAgentPrefVelocity(a0, &rvo.Vector2{X: 1, Y: 1})
	//sim.SetAgentPrefVelocity(a1, &rvo.Vector2{X: -1, Y: 1})
	//sim.SetAgentPrefVelocity(a2, &rvo.Vector2{X: -1, Y: -1})
	//sim.SetAgentPrefVelocity(a3, &rvo.Vector2{X: 1, Y: -1})

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation...\n\n")
}

func main() {
	sim := rvo.NewRVOSimulator(float64(1)/60, 1.5, 5, 1.5, 2, 0.4, 2, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)

	for step := 0; step < 5; step++ {
		sim.DoStep()

		var agentPositions string
		agentPositions = ""
		for j := 0; j < sim.GetNumAgents(); j++ {
			agentPositions = agentPositions + " (" + strconv.FormatFloat(sim.GetAgentPosition(j).X, 'f', 3, 64) + "," + strconv.FormatFloat(sim.GetAgentPosition(j).Y, 'f', 4, 64) + ") "
		}
		fmt.Printf("step=%v  t=%v  %v \n", step, strconv.FormatFloat(sim.GlobalTime, 'f', 3, 64), agentPositions)
	}
}
