// author: asydevc <asydev@163.com>
// date: 2021-02-25

// Package command for build application service.
package service

import (
	"errors"
	"fmt"
	"github.com/asydevc/console/v2/base"
	"github.com/asydevc/console/v2/i"
	"os"
	"regexp"
)

const (
	Description = "Build server file for iris application"
	Name        = "build:service"
)

var (
	regexpFirstLetter = regexp.MustCompile(`^([a-zA-Z0-9])`)
	regexpResetSnake  = regexp.MustCompile(`[_]+([a-zA-Z0-9])`)
	regexpTypeName    = regexp.MustCompile(`^([_a-zA-Z0-9\-]+)`)
)

// Command struct.
type command struct {
	base.Command
	packages map[string]int
}

// New build service instance.
func New() i.ICommand {
	o := &command{packages: make(map[string]int)}
	o.Initialize()
	o.SetDescription(Description)
	o.SetName(Name)
	// service name.
	o.Add(base.NewOption(i.RequiredMode, i.StrValue).SetName("name").SetShortName("n").SetDescription("Service name, no suffix, equal to model name."))
	// application path.
	o.Add(base.NewOption(i.OptionalMode, i.StrValue).SetName("path").SetShortName("p").SetDefaultValue("./app").SetDescription("Application path."))
	// override if file exist.
	//   -o
	//   --override
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetName("override").SetShortName("o").SetDescription("Override if file exist"))
	// with
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetDefaultValue(true).SetName("no-add").SetDescription("Export Add(req *Model) method"))
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetDefaultValue(true).SetName("no-get").SetDescription("Export Get(req *Model) method"))
	o.Add(base.NewOption(i.OptionalMode, i.BoolValue).SetDefaultValue(true).SetName("no-get-by-id").SetDescription("Export GetById(id int) method"))
	// prepared.
	return o
}

// Run command.
func (o *command) Run(console i.IConsole) {
	// variables.
	name := o.GetOption("name").ToString()
	exportName := o.GetOption(name)
	path := o.GetOption("path").ToString() + "/services"
	file := path + "/" + name + ".go"
	// logger.
	console.Info("Command %s: begin.", o.GetName())
	console.Info("        name: %s.", exportName)
	console.Info("        file: %s.", file)
	defer console.Info("Command %s: completed.", o.GetName())
	// file exist for not override.
	if ok, _ := o.fileExist(file); ok && !o.GetOption("override").ToBool() {
		console.PrintError(errors.New(fmt.Sprintf("Command %s: file exist", o.GetName())))
		return
	}
}

func (o *command) dumpHead()          {}
func (o *command) dumpBody()          {}
func (o *command) dumpMethodAdd()     {}
func (o *command) dumpMethodGet()     {}
func (o *command) dumpMethodGetById() {}

// Check model file exist.
func (o *command) fileExist(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
