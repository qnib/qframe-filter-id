package qframe_filter_id

import (
	"C"
	"log"
	"strings"


	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
	"github.com/zpatrick/go-config"
	"github.com/moby/moby/api/types/events"
)

const (
	version = "0.1.0"
	pluginTyp = "filter"
)

type Plugin struct {
	qtypes.Plugin
	types []string
	sendData []string
	sendBack []string
}

func New(qChan qtypes.QChan, cfg config.Config, name string) Plugin {
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, name, version),
		types: []string{},
	}
	p.sendBack = strings.Split(p.CfgStringOr("send-back", ""), ",")
	p.sendData = strings.Split(p.CfgStringOr("send-data", ""), ",")
	return p
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	log.Println("[II] Start id filter '%s'", p.Name)
	myId := qutils.GetGID()
	bg := p.QChan.Data.Join()
	inputs := p.GetInputs()
	for {
		val := bg.Recv()
		switch val.(type) {
		case qtypes.QMsg:
			qm := val.(qtypes.QMsg)
			if qm.SourceID == myId {
				continue
			}
			if len(inputs) != 0 && !qutils.IsInput(inputs, qm.Source) {
				continue
			}
			qm.Type = "filter"
			qm.Source = p.Name
			qm.SourceID = myId
			switch qm.Data.(type) {
			case events.Message:
				if qutils.IsItem(p.sendData, "docker-event") {
					p.QChan.Data.Send(qm)
				}
				if qutils.IsItem(p.sendBack, "docker-event") {
					p.QChan.Back.Send(qm)
				}
			}

		}
	}
}
