package rvosimulator_test

import (
	"testing"

	rvo "../../src/rvosimulator"
)

func TestNewRVOSimulator(t *testing.T) {
	sim := rvo.NewRVOSimulator(1.5, 1, 1, 1, 1, 1, 1, rvo.NewVector2(1, 1))
	// Check if simulator is set right TimeStep
	if sim.TimeStep != 1.5 {
		t.Error("TimeStep: \n実際： ", sim.TimeStep, "\n理想： ", 1.5)
	}
	// Check if simulator is set right GlobalTime
	if sim.GlobalTime != 0.0 {
		t.Error("\n実際： ", sim.GlobalTime, "\n理想： ", 0.0)
	}
	// Check if simulator is set right DefaultAgent
	/*expectedAgent := &rvo.Agent{
		MaxNeighbors:    1,
		MaxSpeed:        1,
		NeighborDist:    1,
		Radius:          1,
		TimeHorizon:     1,
		TimeHorizonObst: 1,
		Velocity:        rvo.NewVector2(1, 1),
	}
	if sim.DefaultAgent != expectedAgent {
		t.Error("\n実際： ", sim.DefaultAgent, "\n理想： ", expectedAgent)
	}*/

	t.Log("TestNewRVOSimulator終了")
}

func TestIsReachedGoal(t *testing.T) {
	sim := rvo.NewRVOSimulator(1.5, 1, 1, 1, 1, 1, 1, rvo.NewVector2(1, 1))
	id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 0, Y: 0})
	sim.SetAgentRadius(id, 1)
	sim.SetAgentGoal(id, &rvo.Vector2{X: 2, Y: 2})

	// Check if agent doesn't reached goal
	if sim.IsReachedGoal() == true {
		t.Error("\n実際： ", sim.IsReachedGoal(), "\n理想： ", false)
	}

	// Check if agent reached goal
	sim.SetAgentPosition(id, rvo.NewVector2(1, 2))
	if sim.IsReachedGoal() == false {
		t.Error("\n実際： ", sim.IsReachedGoal(), "\n理想： ", true)
	}
}

func TestGetGoalVector(t *testing.T) {
	sim := rvo.NewRVOSimulator(1.5, 1, 1, 1, 1, 1, 1, rvo.NewVector2(1, 1))
	id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 0, Y: 0})
	sim.SetAgentRadius(id, 1)
	sim.SetAgentGoal(id, &rvo.Vector2{X: 2, Y: 2})
	// Check if goal vector is right
	expectedVector := &rvo.Vector2{X: 2, Y: 2}
	if sim.GetAgentGoalVector(id) == expectedVector {
		t.Error("\n実際： ", sim.GetAgentGoalVector(id), "\n理想： ", expectedVector)
	}
}
