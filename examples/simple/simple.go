package main

import (
	"fmt"
	"strconv"
	//"math/rand"

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	monitor "github.com/RuiHirano/rvo2-go/monitor"
)


func setupScenario(sim *rvo.RVOSimulator) {

	/*agentNum := 1
	for i := 0; i < agentNum; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 1.5, Y: 0})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: -2, Y: 1})
		id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: 0, Y: -0.5})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: -1, Y: 0})
		id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: 0, Y: 1.5})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 0, Y: -1})
		id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: 0, Y: 0.5})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 0, Y: 1})
		id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: 0, Y: 0.5})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: -1, Y: 0})
		id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: -1.5, Y: 0})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 1, Y: 0})
		id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: -3, Y: 0.1})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 1, Y: 0})
	}*/

	// 6まで同じ。7から値が違う(壁があってもなくても)
	// 9まで壁に止まる。10からすり抜ける
	for i := 0; i < 8; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 2.4, Y: 0 + float64(i)*0.01})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: -2, Y: 0})
	}

	//id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 2.4, Y: 0})
	//	sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: -2, Y: 0})

		//id, _ = sim.AddDefaultAgent(&rvo.Vector2{X: 2.3, Y: 0})
		//sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 2, Y: 0})

	/*for i := 0; i < 20; i++ {
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 1.7, Y: 0 + float64(i)*0.005})
		sim.SetAgentPrefVelocity(id, &rvo.Vector2{X: 2, Y: 0})
	}*/

	/*obstacle1 := []*rvo.Vector2{
		&rvo.Vector2{X: 1, Y: 1},
		&rvo.Vector2{X: -1, Y: 1},
		&rvo.Vector2{X: -1, Y: -1},
	}
	sim.AddObstacle(obstacle1)
	obstacle2 := []*rvo.Vector2{
		&rvo.Vector2{X: -2, Y: -2},
		&rvo.Vector2{X: 2, Y: -2},
		&rvo.Vector2{X: 2, Y: 2},
		&rvo.Vector2{X: -2, Y: 2},
	}
	sim.AddObstacle(obstacle2)*/
	obstacle3 := []*rvo.Vector2{
		&rvo.Vector2{X: 2.2, Y: -0.10},
		&rvo.Vector2{X: 2.2, Y: 0.3},
	}
	sim.AddObstacle(obstacle3)


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
		fmt.Printf("step=%v  t=%v  %v \n\n\n\n", step+1, strconv.FormatFloat(sim.GlobalTime, 'f', 3, 64), agentPositions)

}

func main() {
	timeStep := 0.20
	neighborDist := 0.05 //0.05 // どのくらいの距離の相手をNeighborと認識するか.小さいと衝突する可能性があり、安全ではない
	maxneighbors := 5   // 周り何体を計算対象とするか
	timeHorizon := 0.5
	timeHorizonObst := 0.5
	radius := 0.01  // エージェントの半径
	maxSpeed := 0.1 // エージェントの最大スピード
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
