// author: asydevc <asydev@163.com>
// date: 2021-02-25

package console

import (
	"github.com/asydevc/console/v2/base"
	"github.com/asydevc/console/v2/i"
	"github.com/asydevc/console/v2/s/build/model"
	"github.com/asydevc/console/v2/s/build/service"
	"github.com/asydevc/console/v2/s/docs"
	"github.com/asydevc/console/v2/s/help"
)

// Return default console.
func Default() i.IConsole {
	c := New()
	c.Add(docs.New())
	c.Add(model.New())
	c.Add(service.New())
	return c
}

// Return new console.
func New() i.IConsole {
	o := base.NewConsole()
	o.Add(help.New())
	return o
}
