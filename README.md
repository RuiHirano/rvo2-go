# Optimal Reciprocal Collision Avoidance for Golang

New updates are released
there. There are no explicit version numbers -- all commits on the master
branch are supposed to be stable.

## Building & installing

```
git clone https://github.com/RuiHirano/rvo2-go.git
cd rvo2-go/examples/simple
go build simple.go
./simple
```

## Example code

```go
package main

import (
	"fmt"
	"strconv"

	rvo "../../src/rvosimulator"
)

var (
	goals []*rvo.Vector2
)

func init() {
	goals = make([]*rvo.Vector2, 0)
}

func setupScenario(sim *rvo.RVOSimulator) {

	a0, _ := sim.AddAgentPosition(&rvo.Vector2{X: 0, Y: 0})
	a1, _ := sim.AddAgentPosition(&rvo.Vector2{X: 1, Y: 0})
	a2, _ := sim.AddAgentPosition(&rvo.Vector2{X: 1, Y: 1})
	a3, _ := sim.AddAgentPosition(&rvo.Vector2{X: 0, Y: 1})

	obstaclesPosition := []*rvo.Vector2{
		&rvo.Vector2{X: 0.1, Y: 0.1},
		&rvo.Vector2{X: -0.1, Y: 0.1},
		&rvo.Vector2{X: -0.1, Y: -0.1},
	}
	obstaclesPosition2 := []*rvo.Vector2{
		&rvo.Vector2{X: 0.2, Y: 0.2},
		&rvo.Vector2{X: -0.2, Y: 0.2},
	}
	sim.AddObstacle(obstaclesPosition)
	sim.AddObstacle(obstaclesPosition2)
	sim.ProcessObstacles()

	sim.SetAgentPrefVelocity(a0, &rvo.Vector2{X: 1, Y: 1})
	sim.SetAgentPrefVelocity(a1, &rvo.Vector2{X: -1, Y: 1})
	sim.SetAgentPrefVelocity(a2, &rvo.Vector2{X: -1, Y: -1})
	sim.SetAgentPrefVelocity(a3, &rvo.Vector2{X: 1, Y: -1})

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation...\n\n")
}

func main() {
	sim := rvo.NewRVOSimulator(float64(1)/60, 1.5, 5, 1.5, 2, 0.4, 2, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)

	for step := 0; step < 200; step++ {
		sim.DoStep()

		var agentPositions string
		agentPositions = ""
		for j := 0; j < sim.GetNumAgents(); j++ {
			agentPositions = agentPositions + " (" + strconv.FormatFloat(sim.GetAgentPosition(j).X, 'f', 3, 64) + "," + strconv.FormatFloat(sim.GetAgentPosition(j).Y, 'f', 4, 64) + ") "
		}
		fmt.Printf("step=%v  t=%v  %v \n", step, strconv.FormatFloat(sim.GlobalTime, 'f', 3, 64), agentPositions)
	}
}

```

# Optimal Reciprocal Collision Avoidance

<http://gamma.cs.unc.edu/RVO2/>

[![Build Status](https://travis-ci.org/snape/RVO2.svg?branch=master)](https://travis-ci.org/snape/RVO2)
[![Build status](https://ci.appveyor.com/api/projects/status/0nyp7y4di8x1gh9o/branch/master?svg=true)](https://ci.appveyor.com/project/snape/rvo2)

Copyright 2008 University of North Carolina at Chapel Hill

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Please send all bug reports for the Python wrapper to
[Python-RVO2](https://github.com/sybrenstuvel/Python-RVO2), and bug
report for the RVO2 library itself to [geom@cs.unc.edu](mailto:geom@cs.unc.edu).

The RVO2 authors may be contacted via:

Jur van den Berg, Stephen J. Guy, Jamie Snape, Ming C. Lin, and Dinesh Manocha  
Dept. of Computer Science  
201 S. Columbia St.  
Frederick P. Brooks, Jr. Computer Science Bldg.  
Chapel Hill, N.C. 27599-3175  
United States of America
