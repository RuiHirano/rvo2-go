package main

import (
	"fmt"
	"strconv"
	"math/rand"
	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
	monitor "github.com/RuiHirano/rvo2-go/monitor"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"io/ioutil"
	"log"
)

var (
	fcs *geojson.FeatureCollection
	geofile string
)

func init(){
	geofile = "higashiyama.geojson"
	fcs = loadGeoJson(geofile)
}

func loadGeoJson(fname string) *geojson.FeatureCollection{

	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Print("Can't read file:", err)
		panic("load json")
	}
	fc, _ := geojson.UnmarshalFeatureCollection(bytes)

	return fc
}

func setupScenario(sim *rvo.RVOSimulator) {

	for i := 0; i < 1; i++ {

		id, _ := sim.AddDefaultAgent(&rvo.Vector2{X: 136.9780 + 0.0001 * rand.Float64(), Y:  35.1560 + 0.0001 * rand.Float64()})
		goal := &rvo.Vector2{X: 136.9790 + 0.0001 * rand.Float64(), Y:  35.1450 + 0.0001 * rand.Float64()}
		/*random := rand.Float64()
		fmt.Printf("random: %v", random)
		if random > 0.75{
			goal = &rvo.Vector2{X: -12 + 1*rand.Float64(), Y:  12 + 1*rand.Float64()}
		}else if random > 0.5{
			goal = &rvo.Vector2{X: 12 + 1*rand.Float64(), Y:  -12 + 1*rand.Float64()}
		}else if random > 0.25  {
			goal = &rvo.Vector2{X: -12 + 1*rand.Float64(), Y:  -12 + 1*rand.Float64()}
		}else  {
			goal = &rvo.Vector2{X: 12 + 1*rand.Float64(), Y:  0 + 1*rand.Float64()}
		}*/
		sim.SetAgentGoal(id, goal)
		goalVector := sim.GetAgentGoalVector(id)
		sim.SetAgentPrefVelocity(id, goalVector)
	}

	// Set Obstacle
	for i, feature := range fcs.Features {
		multiPosition := feature.Geometry.(orb.MultiLineString)[0]
		//fmt.Printf("geometry: ", multiPosition)
				rvoObstacle := []*rvo.Vector2{}

		log.Printf("obst: %v\n",i )
		for _, positionArray := range multiPosition{
				position := &rvo.Vector2{
					X: positionArray[0],
					Y: positionArray[1],
				}
				log.Printf("position: %v\n", position)

				rvoObstacle = append(rvoObstacle, position)
		}
		sim.AddObstacle(rvoObstacle)
	}

	sim.ProcessObstacles()

	fmt.Printf("Simulation has %v agents and %v obstacle vertices in it.\n", sim.GetNumAgents(), sim.GetNumObstacleVertices())
	fmt.Printf("Running Simulation...\n\n")
}

func showStatus(sim *rvo.RVOSimulator, step int){
	var agentPositions string
		agentPositions = ""
		for j := 0; j < sim.GetNumAgents(); j++ {
			agentPositions = agentPositions + " (" + strconv.FormatFloat(sim.GetAgentPosition(j).X, 'f', 6, 64) + "," + strconv.FormatFloat(sim.GetAgentPosition(j).Y, 'f', 6, 64) + ") "
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
	neighborDist := 0.00005 // どのくらいの距離の相手をNeighborと認識するか?Neighborとの距離をどのくらいに保つか？ぶつかったと認識する距離？
	maxneighbors := 10   // 周り何体を計算対象とするか
	timeHorizon := 1.0
	timeHorizonObst := 1.0
	radius := 0.00003  // エージェントの半径
	maxSpeed := 0.00001 // エージェントの最大スピード
	sim := rvo.NewRVOSimulator(timeStep, neighborDist, maxneighbors, timeHorizon, timeHorizonObst, radius, maxSpeed, &rvo.Vector2{X: 0, Y: 0})
	setupScenario(sim)
	// monitor 
	mo := monitor.NewMonitor(sim)

	for step := 0; step < 150; step++ {
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
