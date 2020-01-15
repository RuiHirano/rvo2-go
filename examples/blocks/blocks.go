package main

import (
	"fmt"
	"strconv"
	"math/rand"
	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	monitor "github.com/RuiHirano/rvo2-go/monitor"
)

var (
	RAND_MAX int
)

func init() {
	RAND_MAX = 32767
}

func setupScenario(sim *rvo.RVOSimulator) {

	for i := 0; i < 1; i++ {
		for j := 0; j < 1; j++ {
			id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 3 + 0.1*rand.Float64(), Y: 0.8*rand.Float64()})
			sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 0+ 0.01*rand.Float64(), Y: -2 - 0.01*rand.Float64()})
			id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: -5 + 0.1*rand.Float64(), Y: -0.8*rand.Float64()})
			sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 2+ 0.01*rand.Float64(), Y: 1 + 0.01*rand.Float64()})
			id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: -2 + 0.1*rand.Float64(), Y: 5.5 + -0.8*rand.Float64()})
			sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 0+ 0.01*rand.Float64(), Y: -2 + 0.01*rand.Float64()})
			id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: -2 + 0.1*rand.Float64(), Y: -5.5 + -0.8*rand.Float64()})
			sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 0+ 0.01*rand.Float64(), Y: 2 + 0.01*rand.Float64()})
			id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: -2 + 0.1*rand.Float64(), Y: -5.5 + -0.8*rand.Float64()})
			sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 3+ 0.01*rand.Float64(), Y: -2 + 0.01*rand.Float64()})
		}
	}

	// Add (polygonal) obstacles, specifying their vertices in counterclockwise order
	obstacle1 := []*rvo.Vector2{
		&rvo.Vector2{X: -1, Y: 4.0},
		&rvo.Vector2{X: -4.0, Y: 4.0},
		&rvo.Vector2{X: -4.0, Y: 1},
		&rvo.Vector2{X: -1, Y: 1},
	}
	obstacle2 := []*rvo.Vector2{
		&rvo.Vector2{X: 1, Y: 4.0},
		&rvo.Vector2{X: 4.0, Y: 4.0},
		&rvo.Vector2{X: 4.0, Y: 1},
		&rvo.Vector2{X: 1, Y: 1},
	}
	obstacle3 := []*rvo.Vector2{
		&rvo.Vector2{X: -1, Y: -4.0},
		&rvo.Vector2{X: -1, Y: -1},
		&rvo.Vector2{X: -4.0, Y: -1},
		&rvo.Vector2{X: -4.0, Y: -4.0},
	}
	obstacle4 := []*rvo.Vector2{
		&rvo.Vector2{X: 1, Y: -4.0},
		&rvo.Vector2{X: 4.0, Y: -4.0},
		&rvo.Vector2{X: 4.0, Y: -1},
		&rvo.Vector2{X: 1, Y: -1},
	}

	// clockwise order: a wall is formed from inside to outside
	obstacle5 := []*rvo.Vector2{
		&rvo.Vector2{X: 6, Y: 6},
		&rvo.Vector2{X: 6, Y: -6},
		&rvo.Vector2{X: -6, Y: -6},
		&rvo.Vector2{X: -6, Y: 6},
	}
	sim.AddObstacle(obstacle1)
	sim.AddObstacle(obstacle2)
	sim.AddObstacle(obstacle3)
	sim.AddObstacle(obstacle4)
	sim.AddObstacle(obstacle5)
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
	neighborDist := 0.5 
	maxneighbors := 20  
	timeHorizon := 0.5
	timeHorizonObst := 0.5
	radius := 0.01  
	maxSpeed := 0.5 
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)
	// monitor 
	mo := monitor.NewMonitor(sim)

	for step := 0; step < 34; step++ {

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
