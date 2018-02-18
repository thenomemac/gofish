package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/thenomemac/gofish/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
