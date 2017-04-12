package qframe_filter_id

import (
	"C"
	"log"
	"strings"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
	"github.com/zpatrick/go-config"
)

const (
	version = "0.1.0"
)

type Plugin struct {
	qtypes.Plugin
}

func New(qChan qtypes.QChan, cfg config.Config, name string) Plugin {
	p := Plugin{
		Plugin: qtypes.Plugin{
			QChan: qChan,
			Cfg:   cfg,
		},
	}
	p.Version = version
	p.Name = name
	return p
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	log.Println("[II] Start id filter '%s'", p.Name)
	myId := qutils.GetGID()
	bg := p.QChan.Data.Join()
	for {
		val := bg.Recv()
		switch val.(type) {
		case qtypes.QMsg:
			qm := val.(qtypes.QMsg)
			if qm.SourceID == myId {
				continue
			}
			qm.Type = "filter"
			qm.Source = strings.Join(append(strings.Split(qm.Source, "->"), "id"), "->")
			qm.SourceID = myId
			p.QChan.Data.Send(qm)
		}
	}
}
