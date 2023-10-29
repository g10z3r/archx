package plugin

import "github.com/g10z3r/archx/internal/domain/service/anthill/common"

type JoinerPlugin struct {
	next      common.Plugin
	mergeFunc func(chan interface{}) interface{}
}

func (j *JoinerPlugin) IsTerminal() bool {
	return j.next != nil
}

func (e *JoinerPlugin) Next() common.Plugin {
	return e.next
}

func (s *JoinerPlugin) SetNext(p common.Plugin) {
	s.next = p
}

func (j *JoinerPlugin) Execute(input interface{}) interface{} {
	mergedData := j.mergeFunc(input.(chan interface{}))
	return mergedData
}
