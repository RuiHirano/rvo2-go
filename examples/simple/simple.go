package main

import (
	"fmt"
	"strconv"
	"math/rand"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	monitor "github.com/RuiHirano/rvo2-go/monitor"
)


func setupScenario(sim *rvo.RVOSimulator) {

	agentNum := 10
	for i := 0; i < agentNum; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 2 + rand.Float64(), Y: 0 + rand.Float64()})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 2 + rand.Float64(), Y: 1 + rand.Float64()})
	}

	for i := 0; i < agentNum; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 0 + rand.Float64(), Y: 1 + rand.Float64()})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 1 + rand.Float64(), Y: 3 + rand.Float64()})
	}

	obstacle1 := []*rvo.Vector2{
		&rvo.Vector2{X: 0.1, Y: 0.1},
		&rvo.Vector2{X: -0.1, Y: 0.1},
		&rvo.Vector2{X: -0.1, Y: -0.1},
	}
	sim.AddObstacle(obstacle1)
	sim.ProcessObstacles()

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation...\n\n")
}

func main() {
	sim := rvo.NewRVOSimulator(float64(1)/60, 1.5, 5, 1.5, 2, 0.4, 2, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)
	// monitor 
	mo := monitor.NewMonitor()

	for step := 0; step < 50; step++ {
		sim.DoStep()

		var agentPositions string
		agentPositions = ""
		for j := 0; j < sim.GetNumAgents(); j++ {
			agentPositions = agentPositions + " (" + strconv.FormatFloat(sim.GetAgentPosition(j).X, 'f', 3, 64) + "," + strconv.FormatFloat(sim.GetAgentPosition(j).Y, 'f', 4, 64) + ") "
		}
		fmt.Printf("step=%v  t=%v  %v \n", step, strconv.FormatFloat(sim.GlobalTime, 'f', 3, 64), agentPositions)

		mo.AddData(sim)
	}

	// if you want to watch monitor
	err := mo.RunServer()
	if err != nil{
		fmt.Printf("error occor...: ", err)
	}

}
