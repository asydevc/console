// author: asydevc <asydev@163.com>
// date: 2021-02-27

package deregister

import (
	"github.com/asydevc/console/v2/base"
	"github.com/asydevc/console/v2/i"
	"github.com/asydevc/console/v2/s/consul"
)

const (
	Description = "Remove registered consul service"
	Name        = "service:deregister"
)

// Command struct.
type command struct {
	base.Command
}

func New() i.ICommand {
	// normal.
	o := &command{}
	o.Initialize()
	o.SetDescription(Description)
	o.SetName(Name)
	// prepared.
	return o
}

// Run download.
func (o *command) Run(console i.IConsole) {
	consul.Manager(console, o).Deregister()
}
