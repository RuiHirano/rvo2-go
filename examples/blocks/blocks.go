package main

import (
	"fmt"
	"math"
	"math/rand"

	rvo "../../src/rvosimulator"
)

var (
	RAND_MAX int
)

func init() {
	RAND_MAX = 32767
}

func setupScenario(sim *rvo.RVOSimulator) {

	sim.SetTimeStep(0.25)

	sim.SetAgentDefaults(15.0, 10, 5.0, 5.0, 2.0, 2.0, &rvo.Vector2{})

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 55.0 + float64(i)*10.0, Y: 55.0 + float64(j)*10.0})
			sim.SetAgentGoal(id, &rvo.Vector2{X: 75.0, Y: 75.0})
		}
	}

	obstacle1 := []*rvo.Vector2{
		&rvo.Vector2{X: -10.0, Y: -40.0},
		&rvo.Vector2{X: -40.0, Y: -40.0},
		&rvo.Vector2{X: -40.0, Y: -10.0},
		&rvo.Vector2{X: -10.0, Y: -10.0},
	}
	obstacle2 := []*rvo.Vector2{
		&rvo.Vector2{X: 10.0, Y: 40.0},
		&rvo.Vector2{X: 10.0, Y: 10.0},
		&rvo.Vector2{X: 40.0, Y: 10.0},
		&rvo.Vector2{X: 40.0, Y: 40.0},
	}
	obstacle3 := []*rvo.Vector2{
		&rvo.Vector2{X: 10.0, Y: -40.0},
		&rvo.Vector2{X: 40.0, Y: -40.0},
		&rvo.Vector2{X: 40.0, Y: -10.0},
		&rvo.Vector2{X: 10.0, Y: -10.0},
	}
	obstacle4 := []*rvo.Vector2{
		&rvo.Vector2{X: -10.0, Y: -40.0},
		&rvo.Vector2{X: -10.0, Y: -10.0},
		&rvo.Vector2{X: -40.0, Y: -10.0},
		&rvo.Vector2{X: -40.0, Y: -40.0},
	}
	sim.AddObstacle(obstacle1)
	sim.AddObstacle(obstacle2)
	sim.AddObstacle(obstacle3)
	sim.AddObstacle(obstacle4)
	sim.ProcessObstacles()

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation...\n\n")
}

func updateVisualization(sim *rvo.RVOSimulator) {
	fmt.Printf("Time: %v\n", sim.GetGlobalTime())

	for i := 0; i < sim.GetNumAgents(); i++ {
		fmt.Printf("ID: %v,  Position: %v\n", i, sim.GetAgentPosition(i))
	}
}

func setPreferredVelocities(sim *rvo.RVOSimulator) {
	for i := 0; i < sim.GetNumAgents(); i++ {
		goalVector := sim.GetAgentGoalVector(i)

		if rvo.Sqr(goalVector) > 1.0 {
			goalVector = rvo.Normalize(goalVector)
		}

		sim.SetAgentPrefVelocity(i, goalVector)

		/*
		 * Perturb a little to avoid deadlocks due to perfect symmetry.
		 */
		angle := float64(rand.Intn(RAND_MAX)) * 2.0 * math.Pi / float64(RAND_MAX)
		dist := float64(rand.Intn(RAND_MAX)) * 0.0001 / float64(RAND_MAX)

		sim.SetAgentPrefVelocity(i, rvo.Add(sim.GetAgentPrefVelocity(i), rvo.MulOne(&rvo.Vector2{X: math.Cos(angle), Y: math.Sin(angle)}, dist)))
	}
}

func main() {
	sim := rvo.NewEmptyRVOSimulator()
	setupScenario(sim)

	for {
		fmt.Printf("goal %v\n", rvo.Sqr(rvo.Sub(sim.GetAgentPosition(0), sim.GetAgentGoalVector(0))))
		if sim.IsReachedGoal() {
			break
		}

		updateVisualization(sim)

		setPreferredVelocities(sim)

		sim.DoStep()
	}
}
