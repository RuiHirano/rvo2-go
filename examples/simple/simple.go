package main

import (
	"fmt"
	"strconv"
	"math/rand"
	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	monitor "github.com/RuiHirano/rvo2-go/monitor"
)


func setupScenario(sim *rvo.RVOSimulator) {

	for i := 0; i < 50; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 2.4+ 0.01*rand.Float64(), Y: 0 + 0.01*rand.Float64()})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: -2+ 0.01*rand.Float64(), Y: 0+ 0.01*rand.Float64()})
	}

	obstacle := []*rvo.Vector2{
		&rvo.Vector2{X: 2.2, Y: -0.10},
		&rvo.Vector2{X: 2.2, Y: 0.3},
	}
	sim.AddObstacle(obstacle)
	sim.ProcessObstacles()

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation...\n\n")
}

func showStatus(sim *rvo.RVOSimulator, step int){
	var agentPositions string
		agentPositions = ""
		for j := 0; j < sim.GetNumAgents(); j++ {
			agentPositions = agentPositions + " (" + strconv.FormatFloat(sim.GetAgentPosition(j).X, 'f', 3, 64) + "," + strconv.FormatFloat(sim.GetAgentPosition(j).Y, 'f', 4, 64) + ") "
		}
		fmt.Printf("step=%v  t=%v  %v \n", step+1, strconv.FormatFloat(sim.GlobalTime, 'f', 3, 64), agentPositions)

}

func main() {
	timeStep := 0.20
	neighborDist := 0.05 
	maxneighbors := 5  
	timeHorizon := 0.5
	timeHorizonObst := 0.5
	radius := 0.01  
	maxSpeed := 0.1 
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)
	// monitor 
	mo := monitor.NewMonitor(sim)

	for step := 0; step < 15; step++ {
		sim.DoStep()

		showStatus(sim, step)

		// add data for monitor
		mo.AddData(sim)
	}

	// run monitor server
	err := mo.RunServer()
	if err != nil{
		fmt.Printf("error occor...: ", err)
	}

}
