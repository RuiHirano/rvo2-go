package main

import (
	"fmt"
	"math"
	"strconv"

	//rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	"encoding/json"
	"io/ioutil"
	"log"

	
	rvo "../../src/rvosimulator"
)

var (
	goals []*rvo.Vector2
)

func init() {
	goals = make([]*rvo.Vector2, 0)
}

type Coord struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type AreaCoord struct {
	StartLat float32 `json:"slat"`
	EndLat   float32 `json:"elat"`
	StartLon float32 `json:"slon"`
	EndLon   float32 `json:"elon"`
}

type Map struct {
	Id    uint32    `json:"id"`
	Coord AreaCoord `json:"coord"`
}

type Agent struct {
	Id       uint32 `json:"id"`
	Velocity Coord  `json:"velocity"`
	Coord    Coord  `json:"coord"`
	Goal     Coord  `json:"goal"`
}

func readMapData() []Map {
	// JSONファイル読み込み
	bytes, err := ioutil.ReadFile("map.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var mapData []Map

	if err := json.Unmarshal(bytes, &mapData); err != nil {
		log.Fatal(err)
	}
	return mapData
}

func readAgentData() []Agent {
	bytes, err := ioutil.ReadFile("agent.json")
	if err != nil {
		log.Fatal(err)
	}
	// JSONデコード
	var agents []Agent

	if err := json.Unmarshal(bytes, &agents); err != nil {
		log.Fatal(err)
	}
	return agents
}

func scale(agents []Agent, mapData []Map) []Agent {
	height := float32(math.Abs(float64(mapData[0].Coord.EndLat - mapData[0].Coord.StartLat)))
	minLat := float32(math.Min(float64(mapData[0].Coord.EndLat), float64(mapData[0].Coord.StartLat)))
	width := float32(math.Abs(float64(mapData[0].Coord.EndLon - mapData[0].Coord.StartLon)))
	minLon := float32(math.Min(float64(mapData[0].Coord.EndLon), float64(mapData[0].Coord.StartLon)))

	scaledAgents := make([]Agent, 0)
	for _, agent := range agents {
		scaledLat := (agent.Coord.Lat - minLat) / height
		scaledLon := (agent.Coord.Lon - minLon) / width
		velLat := agent.Velocity.Lat / height
		velLon := agent.Velocity.Lon / width
		goalLat := (agent.Goal.Lat - minLat) / height
		goalLon := (agent.Goal.Lon - minLon) / width

		scaledCoord := Coord{
			Lat: scaledLat,
			Lon: scaledLon,
		}

		scaledVelocity := Coord{
			Lat: velLat,
			Lon: velLon,
		}
		goal := Coord{
			Lat: goalLat,
			Lon: goalLon,
		}

		scaledAgents = append(scaledAgents, Agent{
			Id:       agent.Id,
			Coord:    scaledCoord,
			Velocity: scaledVelocity,
			Goal:     goal,
		})
		fmt.Printf("scaledLat: %v\n", scaledVelocity)
	}

	return scaledAgents
}

func invScale(scaledAgents []Agent, mapData []Map) []Agent {
	height := float32(math.Abs(float64(mapData[0].Coord.EndLat - mapData[0].Coord.StartLat)))
	minLat := float32(math.Min(float64(mapData[0].Coord.EndLat), float64(mapData[0].Coord.StartLat)))
	width := float32(math.Abs(float64(mapData[0].Coord.EndLon - mapData[0].Coord.StartLon)))
	minLon := float32(math.Min(float64(mapData[0].Coord.EndLon), float64(mapData[0].Coord.StartLon)))

	agents := make([]Agent, 0)
	for _, scaledAgent := range scaledAgents {
		lat := scaledAgent.Coord.Lat*height + minLat
		lon := scaledAgent.Coord.Lon*width + minLon
		velLat := scaledAgent.Velocity.Lat * height
		velLon := scaledAgent.Velocity.Lon * width
		goalLat := scaledAgent.Goal.Lat*height + minLat
		goalLon := scaledAgent.Goal.Lon*width + minLon

		coord := Coord{
			Lat: lat,
			Lon: lon,
		}

		velocity := Coord{
			Lat: velLat,
			Lon: velLon,
		}

		goal := Coord{
			Lat: goalLat,
			Lon: goalLon,
		}

		agents = append(agents, Agent{
			Id:       scaledAgent.Id,
			Coord:    coord,
			Velocity: velocity,
			Goal:     goal,
		})
	}

	return agents
}

func setupScenario(sim *rvo.RVOSimulator, scaledAgents []Agent) {

	for _, agent := range scaledAgents {
		position := &rvo.Vector2{X: float64(agent.Coord.Lat), Y: float64(agent.Coord.Lon)}
		velocity := &rvo.Vector2{X: float64(agent.Velocity.Lat), Y: float64(agent.Velocity.Lon)}
		id, _ := sim.AddAgentPosition(position)
		sim.SetAgentPrefVelocity(id, velocity)
		sim.SetAgentMaxSpeed(id, rvo.Abs(velocity))
	}

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation..\n\n")
}

func setPreferredVelocities(sim *rvo.RVOSimulator, scaledAgents []Agent) []Agent {

	nextScaledAgents := make([]Agent, 0)
	for i, agent := range scaledAgents {
		// setPrefVelocity
		goal := &rvo.Vector2{
			X: float64(agent.Goal.Lat),
			Y: float64(agent.Goal.Lon),
		}
		goalVector := rvo.Sub(goal, sim.GetAgentPosition(i))

		if rvo.Sqr(goalVector) > 1 {
			goalVector = rvo.Normalize(goalVector)
		}

		sim.SetAgentPrefVelocity(i, goalVector)

		// nextScledAgents
		scaledPosition := sim.GetAgentPosition(i)
		scaledVelocity := sim.GetAgentPrefVelocity(i)
		scaledAgent := Agent{
			Id: uint32(i),
			Coord: Coord{
				Lat: float32(scaledPosition.X),
				Lon: float32(scaledPosition.Y),
			},
			Goal: scaledAgents[i].Goal,
			Velocity: Coord{
				Lat: float32(scaledVelocity.X),
				Lon: float32(scaledVelocity.Y),
			},
		}

		nextScaledAgents = append(nextScaledAgents, scaledAgent)
	}
	return nextScaledAgents
}

func showProgress(nextAgents []Agent, sim *rvo.RVOSimulator, step int) {
	var agentPositions string
	agentPositions = ""
	for _, agent := range nextAgents {
		agentPositions = agentPositions + " (" + strconv.FormatFloat(float64(agent.Coord.Lat), 'f', 4, 64) + "," + strconv.FormatFloat(float64(agent.Coord.Lon), 'f', 4, 64) + ") "
	}
	fmt.Printf("step=%v  t=%v  %v \n", step, strconv.FormatFloat(sim.GlobalTime, 'f', 3, 64), agentPositions)
}

func main() {
	sim := rvo.NewRVOSimulator(1, 1.5, 15, 1.5, 2, 0.04, 0.01, &rvo.Vector2{X: 0, Y: 0})

	mapData := readMapData()
	agents := readAgentData()

	scaledAgents := scale(agents, mapData)
	fmt.Printf("scaledagents is : %v\n", scaledAgents)
	invScaledAgents := invScale(scaledAgents, mapData)
	fmt.Printf("invAgents is : %v\n", invScaledAgents)

	setupScenario(sim, scaledAgents)

	nextScaledAgents := scaledAgents
	for step := 0; step < 30; step++ {
		sim.DoStep()

		nextScaledAgents = setPreferredVelocities(sim, nextScaledAgents)

		nextAgents := invScale(nextScaledAgents, mapData)

		showProgress(nextAgents, sim, step)

	}
}
