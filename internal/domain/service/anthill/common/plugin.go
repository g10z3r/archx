package common

import "context"

type Plugin interface {
	Execute(ctx context.Context, input interface{}) interface{}
	Next() Plugin
	SetNext(Plugin)
	IsTerminal() bool
}
