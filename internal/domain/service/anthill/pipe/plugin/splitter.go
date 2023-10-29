package plugin

import (
	"context"
	"sync"

	"github.com/g10z3r/archx/internal/domain/service/anthill/common"
)

type SplitterPlugin struct {
	next     common.Plugin
	branches []common.Plugin
}

func (s *SplitterPlugin) IsTerminal() bool {
	return false
}

func (e *SplitterPlugin) Next() common.Plugin {
	return e.next
}

func (s *SplitterPlugin) SetNext(p common.Plugin) {
	s.next = p
}

func (s *SplitterPlugin) Execute(ctx context.Context, input interface{}) interface{} {
	var wg sync.WaitGroup
	dataChannel := make(chan interface{}, len(s.branches))

	for _, plugin := range s.branches {
		wg.Add(1)
		go func(plugin common.Plugin, input interface{}) {
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
