package dao

import m "github.com/jtyers/tmaas-model"

type ThreatModelIDCreator struct{}

func NewThreatModelIDCreator() ThreatModelIDCreator {
	return ThreatModelIDCreator{}
}
func (ThreatModelIDCreator) Create(id string) m.ThreatModelID {
	return m.NewThreatModelID(id)
}
