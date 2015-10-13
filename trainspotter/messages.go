// Message Processing

package main

import (
	"encoding/json"
	"fmt"
)

const KX string = "KX"

// Store in memory current berth state
var berths map[string]*Berth = make(map[string]*Berth)

type Berth struct {
	Headcode *string `json:"headcode"`
}

type Messages []Type

type Type struct {
	CA *Message `json:"CA_MSG"`
	CB *Message `json:"CB_MSG"`
	CC *Message `json:"CC_MSG"`
	CT *Message `json:"CT_MSG"`
}

type Message struct {
	AreaID  string  `json:"area_id omitempty"`
	Descr   *string `json:"descr,omitempty"`
	From    *string `json:"from,omitempty"`
	MsgType *string `json:"msg_type,omitempty"`
	Time    *string `json:"time,omitempty"`
	To      *string `json:"to,omitempty"`
}

func Process(body []byte, h Hub) {
	messages := &Messages{}
	err := json.Unmarshal(body, messages)
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}
	for _, message := range *messages {
		var m *Message

		if message.CA != nil {
			m = message.CA
		}
		if message.CB != nil {
			m = message.CB
		}
		if message.CC != nil {
			m = message.CC
		}

		if m != nil {
			if m.AreaID == KX {
				if m.To != nil {
					berths[*m.To] = &Berth{
						Headcode: m.Descr,
					}
				}
				if m.From != nil {
					berths[*m.From] = nil
				}

				d, _ := json.Marshal(berths)
				h.Broadcast <- d
			}
		}
	}
}
