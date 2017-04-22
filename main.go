package main

import (
	"log"
	"fmt"
	"time"
	"github.com/zpatrick/go-config"
	"github.com/qnib/qframe-types"
	"github.com/qnib/qframe-utils"
	"github.com/qnib/qframe-filter-id/lib"
	"github.com/moby/moby/api/types/events"
)


func Run(qChan qtypes.QChan, cfg config.Config, name string) {
	p := qframe_filter_id.New(qChan, cfg, name)
	p.Run()
}

func main() {
	myId := qutils.GetGID()
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"filter.test.send-back": "docker-event",
		"filter.test.send-data": "docker-event",
		"filter.test.inputs": "docker-events",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	p := qframe_filter_id.New(qChan, *cfg, "test")
	go p.Run()
	time.Sleep(2*time.Second)
	dc := qChan.Data.Join()
	bc := qChan.Back.Join()
	qm := qtypes.NewQMsg("collector", "docker-events")
	qm.SourceID = myId
	qm.Msg = "docker-event"
	qm.Data = events.Message{}
	log.Println("Send message")
	qChan.Data.Send(qm)
	gotData := false
	gotBack := false
	for {
		select {
		case msg := <-dc.Read:
			qm := msg.(qtypes.QMsg)
			if qm.SourceID == myId {
				continue
			}
			fmt.Printf("#### Received message on Data-channel: %s\n", qm.Msg)
			gotData = true
		case msg := <-bc.Read:
			qm := msg.(qtypes.QMsg)
			if qm.SourceID == myId {
				continue
			}
			fmt.Printf("#### Received message on Back-channel: %s\n", qm.Msg)
			gotBack = true
		}
		if gotBack && gotData {
			break
		}
	}
}
