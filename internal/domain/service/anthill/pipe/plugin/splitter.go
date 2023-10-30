package plugin

import (
	"context"
	"sync"
)

type SplitterPlugin struct {
	next     Plugin
	branches []Plugin
}

func (s *SplitterPlugin) IsTerminal() bool {
	return false
}

func (e *SplitterPlugin) Next() Plugin {
	return e.next
}

func (s *SplitterPlugin) SetNext(p Plugin) {
	s.next = p
}

func (s *SplitterPlugin) Execute(ctx context.Context, input interface{}) interface{} {
	var wg sync.WaitGroup
	dataChannel := make(chan interface{}, len(s.branches))

	for _, plugin := range s.branches {
		wg.Add(1)
		go func(plugin Plugin, input interface{}) {
			defer wg.Done()
			dataChannel <- plugin.Execute(ctx, input)
		}(plugin, input)
	}

	go func() {
		wg.Wait()
		close(dataChannel)
	}()

	return dataChannel
}
