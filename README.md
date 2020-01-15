# Optimal Reciprocal Collision Avoidance for Golang

New updates are released there. 
There are no explicit version numbers -- all commits on the master branch are supposed to be stable.

## Latest Update: Simulation Monitor
1. import Monitor
```go
import (
	monitor "github.com/RuiHirano/rvo2-go/monitor"
)
```
2. Create Monitor Instance
```go
mo := monitor.NewMonitor(sim)
```

3. Add Result of Each Step
```go
mo.AddData(sim)
```

4. Run Server
```go
err := mo.RunServer()
if err != nil{
	fmt.Printf("error occor...: ", err)
}
```

You can watch simulation monitor at localhost:8000!

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

	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
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

## Parameter

### Global Parameter
|Paramater  |Type  |Description  |
|---|---|---|
|GlobalTime  |float64  | The Global Time of the simulation. Must be positive.|
|TimeStep  |float64  |The time step of the simulation. Must be positive. |

### Agent Parameter
|Paramater  |Type  |Description  |
|---|---|---|
|MaxNeighbors  |int  | The maximal number of other agents the agent takes into account in the navigation. The larger this number, the longer the running time of the simulation. If the number is too low, the simulation will not be safe.|
|MaxSpeed  |float64  |The maximum speed of the agent. Must be non-negative.|
|NeighborDist  |float64  |The maximal distance (center point to center point) to other agents the agent takes into account in the navigation. The larger this number, the longer the running time of the simulation. If the number is too low, the simulation will not be safe. Must be non-negative.|
|Position  |*RVOSimulator.Vector2  |The current position of the agent.|
|PrefVelocity  |float64  |The current preferred velocity of the agent. This is the velocity the agent would take if no other agents or obstacles were around. The simulator computes an actual velocity for the agent that follows the preferred velocity as closely as possible, but at the same time guarantees collision avoidance.|
|Radius  |float64  |The radius of the agent. Must be non-negative. |
|TimeHorizon |float64  |The minimal amount of time for which the agent's velocities that are computed by the simulation are safe with respect to other agents. The larger this number, the sooner this agent will respond to the presence of other agents, but the less freedom the agent has in choosing its velocities. Must be positive. |
|TimeHorizonObst  |float64  |The minimal amount of time for which the agent's velocities that are computed by the simulation are safe with respect to obstacles. The larger this number, the sooner this agent will respond to the presence of obstacles, but the less freedom the agent has in choosing its velocities. Must be positive.|
|Velocity  |*RVOSimulator.Vector2  |The (current) velocity of the agent.|
|Goal  |*RVOSimulator.Vector2  |The Goal of the agent.|

## RVOSimulator Functions

### Simulator Functions
|Function Name  |Params | Return Type  |Description  |
|---|---|---|---|
|NewEmptyRVOSimulator  |None  |*RVOSimulator | |
|NewRVOSimulator  |()  |*RVOSimulator | |
|DoStep  |None  |None | |
|SetTimeStep  |(timeStep float64 )  |None  | |
|GetGlobalTime  |None  |float64  | |
|GetTimeStep  |None  |float64  | |

### Agent Functions
|Function Name  |Params | Return Type  |Description  |
|---|---|---|---|
|AddDefaultAgent  |(position *RVOSimulator.Vector2)  |(id int, err bool) | |
|AddAgent  |float64  |(id int, err bool) | |
|IsReachedGoal  |None  |bool | |
|IsAgentReachedGoal  |(agentID int)  |bool  | |
|GetAgent  |(agentID int)  |*RVOSimulator.Agent  | |
|GetAgentPosition  |(agentID int)  |*RVOSimulator.Vector2  | |
|GetAgentPrefVelocity  |(agentID int)  |*RVOSimulator.Vector2  | |
|GetAgentRadius  |(agentID int)  |float64  | |
|GetAgentVelocity  |(agentID int)  |*RVOSimulator.Vector2  | |
|GetAgentNumAgents  |None  |int  | |
|GetAgentAgentNeighbor  |(agentID int, neighborID int)  |int  | |
|GetAgentMaxNeighbors  |(agentID int)  |int  | |
|GetAgentMaxSpeed  |(agentID int)  |float64  | |
|GetAgentNeighborDist  |(agentID int)  |float64  | |
|GetAgentNumAgentNeighbors  |(agentID int)  |int  | |
|GetAgentNumObstacleNeighbor  |(agentID int)  |int  | |
|GetAgentNumORCALines  |(agentID int)  |int  | |
|GetAgentTimeHorizon  |(agentID int)  |float64  | |
|GetAgentTimeHorizonObst  |(agentID int)  |float64  | |
|GetAgentORCALine  |(agentID int, lineNo int)  |*RVOSimulator.Line  | |
|GetAgentGoalVector  |(agentID int)  |*RVOSimulator.Vector2  | |
|SetAgentDefaults  |(neighborDist float64, maxNeighbors int, timeHorizon float64, timeHorizonObst float64, radius float64, maxSpeed float64, velocity *RVOSimulator.Vector2 )  |None  | |
|SetAgentMaxNeighbors  |(agentID int, maxNeighbors int )  |None  | |
|SetAgentMaxSpeed  |(agentID int, maxSpeed float64 )  |None  | |
|SetAgentNeighborDist  |(agentID int, neighborDist float64 )  |None  | |
|SetAgentPosition  |(agentID int, position *RVOSimulator.Vector2 )  |None  | |
|SetAgentGoal  |(agentID int, goal *RVOSimulator.Vector2 )  |None  | |
|SetAgentPrefVelocity  |(agentID int, velocity *RVOSimulator.Vector2 )  |None  | |
|SetAgentVelocity  |(agentID int, velocity *RVOSimulator.Vector2 )  |None  | |
|SetAgentRadius  |(agentID int, radius float64 )  |None  | |
|SetAgentTimeHorizon  |(agentID int, timeHorizon float64 )  |None  | |
|SetAgentTimeHorizonObst  |(agentID int, timeHorizonObst float64 )  |None  | |

### Obstacle Functions
|Function Name  |Params | Return Type  |Description  |
|---|---|---|---|
|AddObstacle  |(vertices []*RVOSimulator.Vector2)  |(id int, err bool) | |
|ProcessObstacle  |None  |None | |

# Optimal Reciprocal Collision Avoidance

<http://gamma.cs.unc.edu/RVO2/>

this library (rvo2-go) is based on rvo2-library (https://github.com/snape/RVO2).

Please send all bug reports about rvo2-go to issues of
[rvo2-go](https://github.com/RuiHirano/rvo2-go), and bug
report for the RVO2 library itself to [geom@cs.unc.edu](mailto:geom@cs.unc.edu).

The RVO2 authors may be contacted via:

Jur van den Berg, Stephen J. Guy, Jamie Snape, Ming C. Lin, and Dinesh Manocha  
Dept. of Computer Science  
201 S. Columbia St.  
Frederick P. Brooks, Jr. Computer Science Bldg.  
Chapel Hill, N.C. 27599-3175  
United States of America
