package main

import (
	"fmt"
	"strings"
)

type Client struct {
	id      uint32
	name    string
	address string
	phone   string
}

func (self Client) compare(other Content) bool {
	return self.id == other.(*Client).id
}

func (self Client) compareString(other_name string) bool {
	return other_name == self.name
}

func (self *Client) getId() uint32 {
	return self.id
}

func (self *Client) String() string {
	return self.name
}

func (self *Client) load(rstring string) error {
	var client_values []string = strings.Split(rstring, "|")
	if len(client_values) == 4 {
		self.id = uint32(stringToInt(client_values[0]))
		self.name = client_values[1]
		self.address = client_values[2]
		self.phone = client_values[3]
		return nil
	} else {
		return fmt.Errorf("invalid rstring %v excpected 4 values but %d where found", rstring, len(rstring))
	}
}

func (self Client) toString() string {
	return self.String()
}

func (self Client) toRstring() string {
	return fmt.Sprintf("%d|%s|%s|%s", self.id, self.name, self.address, self.phone)
}

func (self *Client) toPartialJson() string {
	return fmt.Sprintf("{\"id\": %d, \"username\": \"%s\", \"name\": \"%s\"}", self.id, self.name, self.phone)
}

func (self *Client) toJson() string {
	return fmt.Sprintf("{\"id\": %d, \"name\": \"%s\",  \"address\": \"%s\", \"phone\": \"%s\"}", self.id, self.name, self.address, self.phone)
}
