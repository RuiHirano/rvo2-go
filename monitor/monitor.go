package monitor

import (
	"log"
	"fmt"
	"sync"
	"net/http"
	"os"
	"path/filepath"
	//"math"
	gosocketio "github.com/mtfelian/golang-socketio"
	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
)

var (
	assetsDir  http.FileSystem
	mu sync.Mutex
)

type StepData struct {
	Agents []*rvo.Agent
	Obstacles [][]*rvo.Vector2
} 

type RVOParam struct {
	TimeStep float64
	NeighborDist float64
	MaxNeighbors int
	TimeHorizon float64
	TimeHorizonObst float64
	Radius float64
	MaxSpeed float64
} 

type Monitor struct {
	Data []*StepData
	RVOParam *RVOParam
}

func NewMonitor(sim *rvo.RVOSimulator) *Monitor{
	param := &RVOParam{
		TimeStep: sim.TimeStep,
		NeighborDist: sim.DefaultAgent.NeighborDist,
		MaxNeighbors: sim.DefaultAgent.MaxNeighbors,
		TimeHorizon: sim.DefaultAgent.TimeHorizon,
		TimeHorizonObst: sim.DefaultAgent.TimeHorizonObst,
		Radius: sim.DefaultAgent.Radius,
		MaxSpeed: sim.DefaultAgent.MaxSpeed,
	}
	m := &Monitor{
		Data: make([]*StepData, 0),
		RVOParam: param,
	}
	// add initial data
	m.AddData(sim)

	return m
}

// assetsFileHandler for static Data
func assetsFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return
	}

	file := r.URL.Path
	//	log.Printf("Open File '%s'",file)
	if file == "/" {
		file = "/index.html"
	}
	f, err := assetsDir.Open(file)
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Printf("can't open file %s: %v\n", file, err)
		return
	}
	http.ServeContent(w, r, file, fi.ModTime(), f)
}


func (d *StepData) GetJson() string {
	jsonAgents := "["
	for i, agent := range d.Agents{

		// orcaLines
		jsonOrcaLines := "["
		for j, line := range agent.OrcaLines{
				//log.Printf("-----Line %v %v\n", line.Point, line.Direction)
			jsonOrcaLine := ""
			/*if j == len(obstacle)-1 {
				// last
				// 図形を閉じるため最後に一つ追加する
				jsonOrcaLine = fmt.Sprintf(`{"x":%f, "y":%f},{"x":%f, "y":%f}`,
				position.X, position.Y, obstacle[0].X, obstacle[0].Y)

			}else{*/
				if j == len(agent.OrcaLines)-1 {
					// last
					// 図形を閉じるため最後に一つ追加する
					jsonOrcaLine = jsonOrcaLine + fmt.Sprintf(`{"point":{"x":%f, "y":%f},"direction":{"x":%f, "y":%f}}`,
				line.Point.X, line.Point.Y, line.Direction.X, line.Direction.Y)
				}else{
					jsonOrcaLine = fmt.Sprintf(`{"point":{"x":%f, "y":%f},"direction":{"x":%f, "y":%f}},`,
					line.Point.X, line.Point.Y, line.Direction.X, line.Direction.Y)
				}

			
			jsonOrcaLines = jsonOrcaLines + jsonOrcaLine
		}
		jsonOrcaLines = jsonOrcaLines + "]"

		// agents
		jsonAgent := ""
		if i == len(d.Agents)-1 {
			// last
			jsonAgent = fmt.Sprintf(`{"id":%d,"position":{"x":%f, "y":%f},"velocity":{"x":%f,"y":%f},"orcaLines":%s}`,
			agent.ID, agent.Position.X, agent.Position.Y, agent.Velocity.X, agent.Velocity.Y, jsonOrcaLines)
		}else{
			jsonAgent = fmt.Sprintf(`{"id":%d,"position":{"x":%f, "y":%f},"velocity":{"x":%f,"y":%f},"orcaLines":%s},`,
			agent.ID, agent.Position.X, agent.Position.Y, agent.Velocity.X, agent.Velocity.Y, jsonOrcaLines)
		}
		jsonAgents = jsonAgents + jsonAgent
	} 
	jsonAgents = jsonAgents + "]"

	

	jsonObstacles := "["
	for i, obstacle := range d.Obstacles{
		// positions
		jsonPositions := "["
		for j, position := range obstacle{
			jsonPosition := ""
			if j == len(obstacle)-1 {
				// last
				// 図形を閉じるため最後に一つ追加する
				jsonPosition = fmt.Sprintf(`{"x":%f, "y":%f},{"x":%f, "y":%f}`,
				position.X, position.Y, obstacle[0].X, obstacle[0].Y)

			}else{
				jsonPosition = fmt.Sprintf(`{"x":%f, "y":%f},`,
				position.X, position.Y)
			}
			jsonPositions = jsonPositions + jsonPosition
		}
		jsonPositions = jsonPositions + "]"

		// obstacles
		jsonObstacle := ""
		if i == len(d.Obstacles)-1 {
			// last
			jsonObstacle = fmt.Sprintf(`{"id":%d, "positions":%s}`,
			i, jsonPositions)
		}else{
			jsonObstacle = fmt.Sprintf(`{"id":%d, "positions":%s},`,
			i, jsonPositions)
		}
		jsonObstacles = jsonObstacles + jsonObstacle
	} 
	jsonObstacles = jsonObstacles + "]"

	s := fmt.Sprintf(`{"agents":%s,"obstacles":%s}`,
		jsonAgents, jsonObstacles)
	return s
}

func (p *RVOParam) GetJson() string {
	s := fmt.Sprintf(`{"timeStep":%f,"neighborDist":%f,"maxNeighbors":%d,"timeHorizon":%f,"timeHorizonObst":%f,"radius":%f,"maxSpeed":%f}`,
		p.TimeStep, p.NeighborDist, p.MaxNeighbors, p.TimeHorizon, p.TimeHorizonObst, p.Radius, p.MaxSpeed)
	return s
}

func (m *Monitor)RunServer() error {
	currentRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	d := filepath.Join(currentRoot, "..", "..", "monitor", "build")

	assetsDir = http.Dir(d)
	log.Println("AssetDir:", assetsDir)

	server := gosocketio.NewServer()

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected from %s as %s", c.IP(), c.Id())
		
		// send data
		jsonDataArray := make([]string, 0)
		for _, data := range m.Data{
			jsonDataArray = append(jsonDataArray, data.GetJson())
		}
		mu.Lock()
		server.BroadcastToAll("rvo", jsonDataArray)
		mu.Unlock()

		// send param
		mu.Lock()
		server.BroadcastToAll("param", m.RVOParam.GetJson())
		mu.Unlock()
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected from %s as %s", c.IP(), c.Id())
	})


	serveMux := http.NewServeMux()

	serveMux.Handle("/socket.io/", server)
	serveMux.HandleFunc("/", assetsFileHandler)
	log.Println("Serving at localhost:8000...")
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 8000), serveMux); err != nil {
		log.Panic(err)
	}

	return nil
}

func (m *Monitor)shiftOrcaLines(agents []*rvo.Agent){
	// m.Dataが空の場合、スキップ
	if len(m.Data) > 0{
		for i, agent := range agents{
			m.Data[len(m.Data)-1].Agents[i].OrcaLines = agent.OrcaLines
		}
	}
}

func (m *Monitor)AddData(sim *rvo.RVOSimulator){
	// to show in monitor
	agents := make([]*rvo.Agent, 0)
	// to change pointa of each agent
	for i := 0; i < sim.GetNumAgents(); i++ {
		agent := *sim.GetAgent(i)
		agents = append(agents, &agent)
	}

	// orcaLinesを一つずらす
	m.shiftOrcaLines(agents)

	// obstacles
	obstacles := make([][]*rvo.Vector2, 0)
	for i := 0; i < sim.GetNumObstacles(); i++ {
		obstacle := sim.GetObstacle(i)
		obst := make([]*rvo.Vector2, 0)
		for _, obs := range obstacle{
			obst = append(obst, &(*obs))
		}
		obstacles = append(obstacles, obst)
	}
	m.Data = append(m.Data, &StepData{
		Agents: agents,
		Obstacles: obstacles,
	})
}