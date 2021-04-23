package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const SERVER_DATA = "./operational_data"
const USERS_DATA = "./operational_data/users.txt"

type State struct {
	users *List
}

func (self *State) init() {
	self.users = new(List)
}

func (self *State) addUser(user *User) {
	self.users.append(user)
	if err := self.saveStateToLocalStorage(); err != nil {
		panic(err)
	}
}

func (self *State) deleteUserById(id uint32) error {
	var user *User = self.getUserById(id)
	if user == nil {
		return fmt.Errorf("User %d doesnt exists", id)
	}
	self.users.remove(user.toString())
	return nil
}

func (self *State) getNewUserId() int {
	return self.users.length
}

func (self *State) getUserById(target_id uint32) *User {
	var user *User
	var current_node *ListNode = self.users.root
	for current_node != nil {
		user = current_node.NodeContent.(*User)
		if user.id == target_id {
			return user
		}
		current_node = current_node.Next
	}
	return nil
}

func (self *State) getUserByUsername(username string) *User {
	target := self.users.exists(username)
	if target != nil {
		return target.(*User)
	} else {
		return nil
	}
}

func (self *State) getUsersAsJson() string {
	// only gets id, username and name
	var users_json string = "["
	var current_user *User
	current_node := self.users.root
	for current_node != nil {
		current_user = current_node.NodeContent.(*User)
		users_json += current_user.toPartialJson()
		if current_node.Next != nil {
			users_json += ","
		}
		current_node = current_node.Next
	}
	return users_json + "]"
}

func (self *State) loadState() error {
	var err error
	if pathExists(SERVER_DATA) {
		if err = self.loadUsers(); err != nil {
			return err
		}
		return err
	} else {
		return fmt.Errorf("'%s' no such file or directory", SERVER_DATA)
	}
}

func (self *State) loadUsers() error {
	if pathExists(USERS_DATA) {
		f, err := ioutil.ReadFile(USERS_DATA)
		if err != nil {
			return err
		}
		var new_user *User
		for _, u := range strings.Split(string(f), "*") {
			new_user = new(User)
			if new_user.load(u) == nil {
				self.users.append(new_user)
			}
		}
		return err
	} else {
		fmt.Printf("No such file or directory: %s\n", USERS_DATA)
		return fmt.Errorf("An error ocurred while loading users: no such file or directory %s", USERS_DATA)
	}
}

func (self *State) saveStateToLocalStorage() error {
	var err error
	fmt.Println("Saving server state...")
	if pathExists(SERVER_DATA) {
		if err = self.saveUsers(); err != nil {
			return err
		}

	} else {
		err = fmt.Errorf("'%s' no such file or directory", SERVER_DATA)
	}
	return err
}

func (self *State) saveUsers() error {
	var rstring string
	var user *User
	for node := self.users.root; node != nil; node = node.Next {
		user = node.NodeContent.(*User)
		rstring += fmt.Sprintf("%s*", user.toRstring())
	}
	fmt.Println("Users saved!")
	return ioutil.WriteFile(USERS_DATA, []byte(rstring), 0755)
}

func (self *State) usernameExists(username string) bool {
	c := self.users.exists(username)
	return c != nil
}
