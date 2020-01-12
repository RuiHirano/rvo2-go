package monitor

import (
	"log"
	"fmt"
	"sync"
	"net/http"
	"os"
	"path/filepath"
	//"encoding/json"
	gosocketio "github.com/mtfelian/golang-socketio"
	//"github.com/googollee/go-socket.io"
	rvo "github.com/RuiHirano/rvo2-go/src/rvosimulator"
)

var (
	port = 7000
	//ioserv     *gosocketio.Server
	assetsDir  http.FileSystem
	mu sync.Mutex
)

type StepData struct {
	Agents []*rvo.Agent
	Obstacles []*rvo.Obstacle
} 

type Monitor struct {
	Data []*StepData
}

func NewMonitor() *Monitor{
	m := &Monitor{
		Data: make([]*StepData, 0),
	}
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
	//log.Printf("data %v\n", d)
	jsonAgents := "["
	for i, agent := range d.Agents{
		jsonAgent := ""
		if i == len(d.Agents)-1 {
			// last
			jsonAgent = fmt.Sprintf(`{"id":%d,"y":%f, "x":%f}`,
			agent.ID, agent.Position.Y, agent.Position.X)
		}else{
			jsonAgent = fmt.Sprintf(`{"id":%d,"y":%f, "x":%f},`,
			agent.ID, agent.Position.Y, agent.Position.X)
		}
		jsonAgents = jsonAgents + jsonAgent
		//log.Printf("jsonAgent %v\n", jsonAgent)
	} 
	jsonAgents = jsonAgents + "]"

	

	jsonObstacles := "["
	log.Printf("obs: %v\n", d.Obstacles)
	for _, obstacle := range d.Obstacles{
		log.Printf("obsta: %v %v\n", obstacle.ID, obstacle.Point, obstacle.NextObstacle)
	}
	/*for i, jsonObstacle := range d.jsonObstacles{
		jsonObstacle = ""
		if i == len(d.jsonObstacles)-1 {
			// last
			jsonObstacle := fmt.Sprintf(`{"id":%d,"lat":%d, "lon":%d}`,
			d.Obstacles[i].ID, d.Obstacles[i].Position.Latitude, d.Obstacles[i].Position.Longitude)
		}else{
			jsonObstacle := fmt.Sprintf(`{"id":%d,"lat":%d, "lon":%d},`,
			d.Obstacles[i].ID, d.Obstacles[i].Position.Latitude, d.Obstacles[i].Position.Longitude)
		}
		jsonObstacles = jsonObstacles + jsonObstacle
	} */
	jsonObstacles = jsonObstacles + "]"

	s := fmt.Sprintf(`{"agents":%s,"obstacles":%s}`,
		jsonAgents, jsonObstacles)
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
		//c.Emit("rvo", jsonDataArray)
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

func (m *Monitor)AddData(sim *rvo.RVOSimulator){
	// to show in monitor
	agents := make([]*rvo.Agent, 0)
	for i := 0; i < sim.GetNumAgents(); i++ {
		//fmt.Printf("agent: %v %v\n",  sim.GetAgent(i))
		agent := *sim.GetAgent(i)
		agents = append(agents, &agent)
	}
	m.Data = append(m.Data, &StepData{
		Agents: agents,
		Obstacles: sim.Obstacles,
	})
}

func (m *Monitor)AddAgents(sim *rvo.RVOSimulator){
	// to show in monitor
	agents := make([]*rvo.Agent, 0)
	for i := 0; i < sim.GetNumAgents(); i++ {
		//fmt.Printf("agent: %v %v\n",  sim.GetAgent(i))
		agent := *sim.GetAgent(i)
		agents = append(agents, &agent)
	}
	m.Data = append(m.Data, &StepData{
		Agents: agents,
		Obstacles: sim.Obstacles,
	})
}

func (m *Monitor)AddObstacles(sim *rvo.RVOSimulator){
	// to show in monitor
	agents := make([]*rvo.Agent, 0)
	for i := 0; i < sim.GetNumAgents(); i++ {
		//fmt.Printf("agent: %v %v\n",  sim.GetAgent(i))
		agent := *sim.GetAgent(i)
		agents = append(agents, &agent)
	}
	m.Data = append(m.Data, &StepData{
		Agents: agents,
		Obstacles: sim.Obstacles,
	})
}