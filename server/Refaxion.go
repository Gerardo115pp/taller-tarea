package main

import (
	"fmt"
	"strings"
)

type Refaxion struct {
	id    uint32
	name  string
	desc  string
	stock uint
}

func (self *Refaxion) compareString(other string) bool {
	return self.name == other
}

func (self *Refaxion) getId() uint32 {
	return self.id
}

func (self *Refaxion) load(rstring string) error {
	var vehicle_values []string = strings.Split(rstring, "|")
	if len(vehicle_values) == 4 {
		self.id = uint32(stringToInt(vehicle_values[0]))
		self.name = vehicle_values[1]
		self.desc = vehicle_values[2]
		self.stock = uint(stringToInt(vehicle_values[3]))
		return nil
	} else {
		return fmt.Errorf("invalid rstring %v excpected 4 values but %d where found", rstring, len(rstring))
	}
}

func (self *Refaxion) toString() string {
	return self.name
}

func (self *Refaxion) toRstring() string {
	return fmt.Sprintf("%d|%s|%s|%d", self.id, self.name, self.desc, self.stock)
}

func (self *Refaxion) toPartialJson() string {
	return fmt.Sprintf("{\"id\": %d, \"username\": \"%s\", \"name\": \"%d\"}", self.id, fmt.Sprintf("%s(%d)", self.name, self.stock), fmt.Sprintf("%s...", self.desc[:clamp(len(self.desc), 1, 15)]))
}

func (self *Refaxion) toJson() string {
	return fmt.Sprintf("{\"id\": %d, \"name\": \"%s\",  \"description\": \"%s\", \"stock\": %d}", self.id, self.name, self.desc, self.stock)
}
