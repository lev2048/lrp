package status

import (
	"fmt"
	"math"
	"time"

	"github.com/lev2048/agent"
)

type Data struct {
	Cpu           string `json:"cpu"`
	Mem           string `json:"mem"`
	Upload        string `json:"upload"`
	Download      string `json:"download"`
	TotalUpload   string `json:"totalUpload"`
	TotalDownload string `json:"totalDownload"`
}

type Monitor struct {
	exit  chan bool
	Data  Data
	agent *agent.Agent
}

func NewMonitor() *Monitor {
	return &Monitor{
		exit: make(chan bool),
	}
}

func (m *Monitor) Start() {
	m.agent = agent.NewAgent("agent")
	m.agent.Start(false)
	go func(exit chan bool) {
		for {
			select {
			case <-exit:
				return
			default:
				data := m.agent.GetData()
				if data.MemTotal == 0 {
					time.Sleep(time.Duration(1) * time.Second)
					continue
				}
				m.Data = Data{
					fmt.Sprintf("%d%%", int(math.Ceil(data.CPU*100))),
					fmt.Sprintf("%d%%", int(math.Ceil((float64(data.MemUsed)/float64(data.MemTotal))*100))),
					agent.UnitConver(float64(data.NetworkTx)) + "/s",
					agent.UnitConver(float64(data.NetworkRx)) + "/s",
					agent.UnitConver(float64(data.NetworkOut)),
					agent.UnitConver(float64(data.NetworkIn)),
				}
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}(m.exit)
}

func (m *Monitor) Info() interface{} {
	return m.Data
}

func (m *Monitor) Stop() {
	close(m.exit)
	m.agent.Stop()
}
