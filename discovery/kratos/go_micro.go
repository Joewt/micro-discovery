package kratos

import (
	"encoding/json"
)

const prefix = "/micro/registry"

type Service struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []*Endpoint       `json:"endpoints"`
	Nodes     []*Node           `json:"nodes"`
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}

type Endpoint struct {
	Name     string            `json:"name"`
	Request  *Value            `json:"request"`
	Response *Value            `json:"response"`
	Metadata map[string]string `json:"metadata"`
}

type Value struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []*Value `json:"values"`
}

func encode(s *Service) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "",err
	}
	return string(b), nil
}

func decode(ds []byte) (s *Service, err error) {
	err = json.Unmarshal(ds, &s)
	return
}