package nfdc

import (
	"time"

	"github.com/zjkmxy/go-ndn/pkg/log"
	"github.com/zjkmxy/go-ndn/pkg/ndn"
	mgmt "github.com/zjkmxy/go-ndn/pkg/ndn/mgmt_2022"
)

type NfdMgmtCmd struct {
	Module  string
	Cmd     string
	Args    *mgmt.ControlArgs
	Retries int
}

type NfdMgmtThread struct {
	// engine
	engine ndn.Engine
	// channel for management commands
	channel chan NfdMgmtCmd
	// stop the management thread
	stop chan bool
}

func NewNfdMgmtThread(engine ndn.Engine) *NfdMgmtThread {
	return &NfdMgmtThread{
		engine:  engine,
		channel: make(chan NfdMgmtCmd, 4096),
		stop:    make(chan bool),
	}
}

func (m *NfdMgmtThread) Start() {
	for {
		select {
		case cmd := <-m.channel:
			for i := 0; i < cmd.Retries || cmd.Retries < 0; i++ {
				err := m.engine.ExecMgmtCmd(cmd.Module, cmd.Cmd, cmd.Args)
				if err != nil {
					log.Errorf("nfdc %s %s failed: %s %+v [%d]", cmd.Module, cmd.Cmd, cmd.Args.Name, err, i)
					time.Sleep(100 * time.Millisecond)
				} else {
					time.Sleep(1 * time.Millisecond)
					break
				}
			}
		case <-m.stop:
			return
		}
	}
}

func (m *NfdMgmtThread) Stop() {
	m.stop <- true
	close(m.channel)
	close(m.stop)
}

func (m *NfdMgmtThread) Exec(mgmt_cmd NfdMgmtCmd) {
	m.channel <- mgmt_cmd
}
