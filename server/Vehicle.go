package main

import (
	"fmt"
	"strings"
)

type Vehicle struct {
	id      uint32
	plates  string
	brand   string
	model   string
	release string
	client  uint32
}

func (self *Vehicle) compareString(other string) bool {
	return self.plates == other
}

func (self *Vehicle) getId() uint32 {
	return self.id
}

func (self *Vehicle) load(rstring string) error {
	var vehicle_values []string = strings.Split(rstring, "|")
	if len(vehicle_values) == 6 {
		self.id = uint32(stringToInt(vehicle_values[0]))
		self.plates = vehicle_values[1]
		self.brand = vehicle_values[2]
		self.model = vehicle_values[3]
		self.release = vehicle_values[4]
		self.client = uint32(stringToInt(vehicle_values[5]))
		return nil
	} else {
		return fmt.Errorf("invalid rstring %v excpected 4 values but %d where found", rstring, len(rstring))
	}
}

func (self *Vehicle) toString() string {
	return self.plates
}

func (self *Vehicle) toRstring() string {
	return fmt.Sprintf("%d|%s|%s|%s|%s|%d", self.id, self.plates, self.brand, self.model, self.release, self.client)
}

func (self *Vehicle) toPartialJson() string {
	return fmt.Sprintf("{\"id\": %d, \"username\": \"%s\", \"name\": \"%s\"}", self.id, fmt.Sprintf("%s-%s(%s)", self.model, self.brand, self.release), self.plates)
}

func (self *Vehicle) toJson() string {
	return fmt.Sprintf("{\"id\": %d, \"plates\": \"%s\",  \"brand\": \"%s\", \"model\": \"%s\",\"release\": \"%s\",\"client\": \"%d\"}", self.id, self.plates, self.brand, self.model, self.release, self.client)
}
