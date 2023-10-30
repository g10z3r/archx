package plugin

type JoinerPlugin struct {
	next      Plugin
	mergeFunc func(chan interface{}) interface{}
}

func (j *JoinerPlugin) IsTerminal() bool {
	return j.next != nil
}

func (e *JoinerPlugin) Next() Plugin {
	return e.next
}

func (s *JoinerPlugin) SetNext(p Plugin) {
	s.next = p
}

func (j *JoinerPlugin) Execute(input interface{}) interface{} {
	mergedData := j.mergeFunc(input.(chan interface{}))
	return mergedData
}
