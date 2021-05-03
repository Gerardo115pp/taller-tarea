package main

import (
	"fmt"
	"strings"
)

type Service struct {
	id               uint32
	start            string
	end              string
	failure_desc     string
	vehicle          uint32
	refaxions_needed string // most be a json
}

func (self *Service) compareString(other string) bool {
	return self.toString() == other
}

func (self *Service) getId() uint32 {
	return self.id
}

func (self *Service) load(rstring string) error {
	var vehicle_values []string = strings.Split(rstring, "|")
	if len(vehicle_values) == 6 {
		self.id = uint32(stringToInt(vehicle_values[0]))
		self.start = vehicle_values[1]
		self.end = vehicle_values[2]
		self.failure_desc = vehicle_values[3]
		self.vehicle = uint32(stringToInt(vehicle_values[4]))
		self.refaxions_needed = vehicle_values[5]
		return nil
	} else {
		return fmt.Errorf("invalid rstring %v excpected 4 values but %d where found", rstring, len(rstring))
	}
}

func (self *Service) toString() string {
	return fmt.Sprintf("service-%d-%d", self.id, self.vehicle)
}

func (self *Service) toRstring() string {
	return fmt.Sprintf("%d|%s|%s|%s|%d|%s", self.id, self.start, self.end, self.failure_desc, self.vehicle, self.refaxions_needed)
}

func (self *Service) toPartialJson() string {
	return fmt.Sprintf("{\"id\": %d, \"username\": \"%s\", \"name\": \"%s\"}", self.id, self.toString(), fmt.Sprintf("%s - %s", self.start, self.end))
}

func (self *Service) toJson() string {
	return fmt.Sprintf("{\"id\": %d, \"start\": \"%s\",  \"end\": \"%s\", \"failure\": \"%s\",\"vehicle\": \"%d\",\"extras\": %s}", self.id, self.start, self.end, self.failure_desc, self.vehicle, self.refaxions_needed)
}
