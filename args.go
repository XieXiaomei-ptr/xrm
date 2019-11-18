package main

import (
	"flag"
	"strings"
)

type ArgsHanler interface {
	Parse()
}

type StringArgsHanler struct {
	splitValue []string
	value      string
	separator  string
}

func (self *StringArgsHanler) Parse() {
	if len(self.value) == 0 {
		return
	}
	if len(self.separator) > 0 {
		self.splitValue = strings.Split(self.value, self.separator)
	} else {
		self.splitValue = []string{self.value}
	}
}

type BoolArgsHanler struct {
	value bool
}

func (self *BoolArgsHanler) Parse() {
}

type Args struct {
	argsHanlers map[string]ArgsHanler
}

func MakeArgs() *Args {
	result := &Args{}
	result.argsHanlers = make(map[string]ArgsHanler)
	return result
}

func (self *Args) AddStringArg(name string, defaultValue string, usage string, separator string) bool {
	_, ok := self.argsHanlers[name]
	if ok {
		return false
	}
	handler := &StringArgsHanler{}
	handler.separator = separator
	flag.StringVar(&handler.value, name, defaultValue, usage)
	self.argsHanlers[name] = handler
	return true
}

func (self *Args) AddBoolArg(name string, defaultValue bool, usage string) bool {
	_, ok := self.argsHanlers[name]
	if ok {
		return false
	}
	handler := &BoolArgsHanler{}
	flag.BoolVar(&handler.value, name, defaultValue, usage)
	self.argsHanlers[name] = handler
	return true
}

func (self *Args) Parse() bool {
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		return false
	}
	for _, handler := range self.argsHanlers {
		handler.Parse()
	}
	return true
}

func (self *Args) GetStringVar(name string) []string {
	handler, ok := self.argsHanlers[name]
	if !ok {
		return nil
	}
	strHandler := handler.(*StringArgsHanler)
	return strHandler.splitValue
}

func (self *Args) GetBoolVar(name string) bool {
	handler, ok := self.argsHanlers[name]
	if !ok {
		return false
	}
	strHandler := handler.(*BoolArgsHanler)
	return strHandler.value
}
