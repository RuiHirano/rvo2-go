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
		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: float64(136.974694) + 0.0001*rand.Float64(), Y:  35.158200 + 0.0001*rand.Float64()})
		goal := &rvo.Vector2{X: 136.974640, Y:  35.157671}
		sim.SetAgentGoal(id, goal)
		goalVector := sim.GetAgentGoalVector(id)
		sim.SetAgentPrefVelocity(id, goalVector)
	}

	/*obstacle := []*rvo.Vector2{
		&rvo.Vector2{X: 236.121324222, Y: 135.331234123},
		&rvo.Vector2{X: 233.12124123422, Y: 236.32134123433},
	}
	sim.AddObstacle(obstacle)
	sim.ProcessObstacles()*/

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

func setPreferredVelocities(sim *rvo.RVOSimulator) {
	for i := 0; i < sim.GetNumAgents(); i++ {
		goalVector := sim.GetAgentGoalVector(i)
		sim.SetAgentPrefVelocity(i, goalVector)
	}
}

func main() {
	timeStep := float64(1)
	neighborDist := 0.0005 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 10   // 周り何体を計算対象とするか
	timeHorizon := 1.0
	timeHorizonObst := 1.0
	radius := 0.0001  // エージェントの半径
	maxSpeed := 0.001 // エージェントの最大スピード
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)
	// monitor 
	mo := monitor.NewMonitor(sim)

	for step := 0; step < 50; step++ {
		sim.DoStep()
		setPreferredVelocities(sim)

		showStatus(sim, step)

		// add data for monitor
		mo.AddData(sim)
		
		if sim.IsReachedGoal(){
			break
		}
	}

	// run monitor server
	err := mo.RunServer()
	if err != nil{
		fmt.Printf("error occor...: ", err)
	}

}
