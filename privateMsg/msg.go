package main

import (
	"encoding/json"
	"fmt"
)

type Msg struct {
	Fromid  int    `json:"fromid"`
	Toid    int    `json:"toid"`
	Token   string `json:"token,omitempty"`
	Content string `json:"content"`
}

func (m *Msg) Encode() []byte {

	buf, _ := json.Marshal(m)

	return buf
}

func (m *Msg) Decode(data []byte) error {

	err := json.Unmarshal(data, m)
	if err != nil {
		return fmt.Errorf("protocol error:", err.Error())
	}

	return nil
}
