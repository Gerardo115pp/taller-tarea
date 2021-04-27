package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"strings"
)

type UserGroup uint32

const (
	ADMIN UserGroup = 0
	USER  UserGroup = 1
)

type User struct {
	id       uint32
	username string
	name     string
	password uint64
	address  string
	phone    string
	group    UserGroup
}

func (self User) compare(other Content) bool {
	return self.id == other.(*User).id
}

func (self User) compareString(other_username string) bool {
	return self.username == other_username
}

func (self *User) getId() uint32 {
	return self.id
}

func (self *User) getNormalizedName() string {
	return strings.Replace(self.name, "+", " ", 3)
}

func (self *User) load(rstring string) error {
	var err error
	var fields []string = strings.Split(rstring, "|")
	if len(fields) == 7 {
		self.id = uint32(stringToInt(fields[0]))
		self.username = fields[1]
		self.name = fields[2]
		self.password = stringToUint64(fields[3])
		self.address = fields[4]
		self.phone = fields[5]
		self.group = UserGroup(stringToInt(fields[6]))
	} else {
		err = fmt.Errorf("Field length is %d, expected 7", len(fields))
	}
	return err
}

func (self *User) String() string {
	return fmt.Sprintf("%s:%d", self.username, self.group)
}

func (self User) toString() string {
	return self.username
}

func (self *User) toPartialJson() string {
	return fmt.Sprintf("{\"id\": %d, \"username\": \"%s\", \"name\": \"%s\"}", self.id, self.username, self.getNormalizedName())
}

func (self *User) toJson() string {
	return fmt.Sprintf("{\"id\": %d, \"username\": \"%s\", \"name\": \"%s\", \"password\": \"\", \"address\": \"%s\", \"phone\": \"%s\", \"user-group\": %d}", self.id, self.username, self.getNormalizedName(), self.address, self.phone, self.group)
}

func (self *User) toRstring() string {
	return fmt.Sprintf("%d|%s|%s|%d|%s|%s|%d", self.id, self.username, self.name, self.password, self.address, self.phone, self.group)
}

func createUser(id int, user_params ...string) *User {
	/**
	user parameters most be sorted as follows:

	0: username
	1: name
	2: password
	3: address
	4: phone
	5: group
	*/

	var new_user *User = new(User)
	new_user.id = uint32(id)
	new_user.username = user_params[0]
	new_user.name = user_params[1]
	if user_params[2] != "" {
		new_user.password = shaAsInt64(user_params[2])
	} else {
		new_user.password = 0
	}
	new_user.address = user_params[3]
	new_user.phone = user_params[4]
	new_user.group = UserGroup(stringToInt(user_params[5]))
	return new_user
}

func composeName(name string, mname string, lname string) string {
	return fmt.Sprintf("%s+%s+%s", name, mname, lname)
}

func shaAsInt64(s string) uint64 {
	var hash_bytes [20]byte = sha1.Sum([]byte(s))
	return binary.BigEndian.Uint64(hash_bytes[:])
}
