// Copyright (C) 2013-2017, The MetaCurrency Project (Eric Harris-Braun, Arthur Brock, et. al.)
// Use of this source code is governed by GPLv3 found in the LICENSE file
//----------------------------------------------------------------------------------------
// Nucleus provides an interface for an execution environment interface for chains and their entries
// and factory code for creating nucleii instances

package holochain

import (
	"errors"
	"fmt"
	"strings"
)

type NucleusFactory func(h *Holochain, code string) (Nucleus, error)

type InterfaceSchemaType int

const (
	STRING InterfaceSchemaType = iota
	JSON
)

type Interface struct {
	Name   string
	Schema InterfaceSchemaType
}

type Nucleus interface {
	Type() string
	ValidateEntry(def *EntryDef, entry interface{}) error
	expose(iface Interface) error
	Interfaces() (i []Interface)
	Call(iface string, params interface{}) (interface{}, error)
}

var nucleusFactories = make(map[string]NucleusFactory)

// InterfaceSchema returns a functions schema type
func InterfaceSchema(n Nucleus, name string) (InterfaceSchemaType, error) {
	i := n.Interfaces()
	for _, f := range i {
		if f.Name == name {
			return f.Schema, nil
		}
	}
	return -1, errors.New("function not found: " + name)
}

// RegisterNucleus sets up a Nucleus to be used by the CreateNucleus function
func RegisterNucleus(name string, factory NucleusFactory) {
	if factory == nil {
		panic("Nucleus factory for type %s does not exist." + name)
	}
	_, registered := nucleusFactories[name]
	if registered {
		panic("Nucleus factory for type %s already registered. " + name)
	}
	nucleusFactories[name] = factory
}

// RegisterBultinNucleii adds the built in nucleus types to the factory hash
func RegisterBultinNucleii() {
	RegisterNucleus(ZygoNucleusType, NewZygoNucleus)
}

// CreateNucleus returns a new Nucleus of the given type
func CreateNucleus(h *Holochain, schema string, code string) (Nucleus, error) {

	factory, ok := nucleusFactories[schema]
	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		available := make([]string, 0)
		for k, _ := range nucleusFactories {
			available = append(available, k)
		}
		return nil, errors.New(fmt.Sprintf("Invalid nucleus name. Must be one of: %s", strings.Join(available, ", ")))
	}

	return factory(h, code)
}