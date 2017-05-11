package qframe_filter_id

import (
	"C"
	"fmt"
	"strings"

	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
	"github.com/zpatrick/go-config"
)

const (
	version   = "0.2.2"
	pluginTyp = "filter"
	pluginPkg = "id"
)

type Plugin struct {
	qtypes.Plugin
	types    []string
	sendData []string
	sendBack []string
}

func New(qChan qtypes.QChan, cfg config.Config, name string) Plugin {
	p := Plugin{
		Plugin: qtypes.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, name, version),
		types:  []string{},
	}
	p.sendBack = strings.Split(p.CfgStringOr("send-back", ""), ",")
	p.sendData = strings.Split(p.CfgStringOr("send-data", ""), ",")
	return p
}

// Run fetches everything from the Data channel and flushes it to stdout
func (p *Plugin) Run() {
	p.Log("notice", fmt.Sprintf("Start filter v%s", p.Version))
	myId := qutils.GetGID()
	dc := p.QChan.Data.Join()
	inputs := p.GetInputs()

	for {
		select {
		case val := <-dc.Read:
			switch val.(type) {
			case qtypes.QMsg:
				qm := val.(qtypes.QMsg)
				if qm.SourceID == myId {
					continue
				}
				if len(inputs) != 0 && !qutils.IsInput(inputs, qm.Source) {
					continue
				}
				if qutils.IsItem(p.sendData, qm.Source) {
					qm.SourceID = myId
					p.QChan.Data.Send(qm)
				}
				if qutils.IsItem(p.sendBack, qm.Source) {
					p.QChan.Back.Send(qm)
				}
			}
		}
	}
}
