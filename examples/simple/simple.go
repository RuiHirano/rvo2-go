package main

import (
	"fmt"
	"strconv"
	"math/rand"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	monitor "github.com/RuiHirano/rvo2-go/monitor"
)


func setupScenario(sim *rvo.RVOSimulator) {

	agentNum := 100
	for i := 0; i < agentNum; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 2 + rand.Float64(), Y: 0 + rand.Float64()})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: -2 + rand.Float64(), Y: -1 + rand.Float64()})
	}

	for i := 0; i < agentNum; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 0 + rand.Float64(), Y: -1 + rand.Float64()})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 1 + rand.Float64(), Y: 2 + rand.Float64()})
	}

	/*obstacle1 := []*rvo.Vector2{
		&rvo.Vector2{X: 1, Y: 1},
		&rvo.Vector2{X: -1, Y: 1},
		&rvo.Vector2{X: -1, Y: -1},
	}
	sim.AddObstacle(obstacle1)
	obstacle2 := []*rvo.Vector2{
		&rvo.Vector2{X: 2, Y: 2},
		&rvo.Vector2{X: -2, Y: 2},
		&rvo.Vector2{X: -2, Y: -2},
		&rvo.Vector2{X: 2, Y: -2},
	}
	sim.AddObstacle(obstacle2)
	sim.ProcessObstacles()*/

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation...\n\n")
}

func main() {
	timeStep := 0.20
	neighborDist := 0.05 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 100   // 周り何体を計算対象とするか
	timeHorizon := 1.0
	timeHorizonObst := 1.0
	radius := 0.01  // エージェントの半径
	maxSpeed := 0.1 // エージェントの最大スピード
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)
	// monitor 
	mo := monitor.NewMonitor(sim)

	for step := 0; step < 200; step++ {
		sim.DoStep()

		var agentPositions string
		agentPositions = ""
		for j := 0; j < sim.GetNumAgents(); j++ {
			agentPositions = agentPositions + " (" + strconv.FormatFloat(sim.GetAgentPosition(j).X, 'f', 3, 64) + "," + strconv.FormatFloat(sim.GetAgentPosition(j).Y, 'f', 4, 64) + ") "
		}
		fmt.Printf("step=%v  t=%v  %v \n", step, strconv.FormatFloat(sim.GlobalTime, 'f', 3, 64), agentPositions)

		// add data for monitor
		mo.AddData(sim)
	}

	// run monitor server
	err := mo.RunServer()
	if err != nil{
		fmt.Printf("error occor...: ", err)
	}

}
