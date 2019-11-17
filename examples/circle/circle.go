package main

import (
	"fmt"
	"math"

	rvo "../../src/rvosimulator"
)

var (
	goals []*rvo.Vector2
)

func init() {
	goals = make([]*rvo.Vector2, 0)
}

func setupScenario(sim *rvo.RVOSimulator) {
	sim.SetTimeStep(0.25)
	sim.SetAgentDefaults(15.0, 10, 10.0, 10.0, 1.5, 2.0, &rvo.Vector2{}) // where is velocity property ?

	agentNum := 20
	goals = make([]*rvo.Vector2, agentNum)
	for i := 0; i < agentNum; i++ {
		position := &rvo.Vector2{
			X: math.Cos(float64(i) * 2.0 * math.Pi / float64(agentNum)),
			Y: math.Sin(float64(i) * 2.0 * math.Pi / float64(agentNum)),
		}
		_, err := sim.AddAgent1(position)

		if !err {
			t := (i + agentNum/2) % agentNum
			goals[t] = sim.GetAgentPosition(i)
		}
	}
}

func updateVisualization(sim *rvo.RVOSimulator) {
	fmt.Printf("Time: %v\n", sim.GetGlobalTime())

	for i := 0; i < sim.GetNumAgents(); i++ {
		fmt.Printf("ID: %v,  Position: %v\n", i, sim.GetAgentPosition(i))
	}
}

func setPreferredVelocities(sim *rvo.RVOSimulator) {
	for i := 0; i < sim.GetNumAgents(); i++ {
		// goal - agentPosition
		goalVector := rvo.Sub(goals[i], sim.GetAgentPosition(i))

		if rvo.Sqr(goalVector) > 1 {
			goalVector = rvo.Normalize(goalVector)
		}

		sim.SetAgentPrefVelocity(i, goalVector)
	}
}

func reachedGoal(sim *rvo.RVOSimulator) bool {
	for i := 0; i < sim.GetNumAgents(); i++ {
		// ゴールまでの距離の二乗がエージェント半径の距離の二乗よりも大きければまだ到達していない
		if rvo.Sqr(rvo.Sub(sim.GetAgentPosition(i), goals[i])) > sim.GetAgentRadius(i)*sim.GetAgentRadius(i) {
			return false
		}
	}
	return true
}

func main() {
	sim := rvo.NewRVOSimulatorBlank()

	setupScenario(sim)
	for {
		if reachedGoal(sim) {
			fmt.Printf("Goal \n ")
			break
		}

		updateVisualization(sim)

		setPreferredVelocities(sim)
		sim.DoStep()
	}
}
