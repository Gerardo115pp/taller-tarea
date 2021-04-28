package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const SERVER_DATA = "./operational_data"

const USERS_STATE = "users"
const CLIENTS_STATE = "clients"
const VEHICLES_STATE = "vehicles"

type State struct {
	users    *List
	clients  *List
	vehicles *List
}

func (self *State) init() {
	self.users = new(List)
	self.clients = new(List)
	self.vehicles = new(List)
}

func (self *State) addItemToState(item Content, type_name string) {
	var state_slice *List
	state_slice, err := self.getStateSliceByName(type_name)
	if err != nil {
		logFatal(err)
	}
	state_slice.append(item)
	if err := self.saveStateToLocalStorage(); err != nil {
		panic(err)
	}
}

func (self *State) deleteItemById(id uint32, type_name string) error {
	var target Content = self.getItemById(id, type_name)
	var target_slice *List
	target_slice, err := self.getStateSliceByName(type_name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if target == nil {
		return fmt.Errorf("%s %d doesnt exists", type_name, id)
	}
	target_slice.remove(target.toString())
	return nil
}

func (self *State) composeSaveFile(type_name string) string {
	return fmt.Sprintf("%s/%s.txt", SERVER_DATA, type_name)
}

func (self *State) getTypeByName(type_name string) Content {
	switch type_name {
	case USERS_STATE:
		return new(User)
	case CLIENTS_STATE:
		return new(Client)
	case VEHICLES_STATE:
		return new(Vehicle)
	default:
		logFatal(fmt.Errorf("no type for type name: %s", type_name))
		return nil
	}
}

func (self *State) getStateSliceByName(slice_name string) (*List, error) {
	switch slice_name {
	case USERS_STATE:
		return self.users, nil
	case CLIENTS_STATE:
		return self.clients, nil
	case VEHICLES_STATE:
		return self.vehicles, nil
	default:
		return nil, fmt.Errorf("No slice for name %s", slice_name)
	}
}

func (self *State) getItemsAsOptions(type_name string) string {
	var target_slice *List
	var items_options []string
	target_slice, err := self.getStateSliceByName(type_name)
	if err != nil {
		fmt.Printf("Error on getITemsAsOptions: %s\n", err.Error())
		return "[]"
	}
	items_options = target_slice.mapFunc(func(c *ListNode) string {
		return fmt.Sprintf("{\"name\": \"%s\", \"value\": %d}", c.NodeContent.toString(), c.NodeContent.getId())
	})
	return strings.Join(items_options, ",")
}

func (self *State) getNewItemId(type_name string) int {
	item_state, err := self.getStateSliceByName(type_name)
	if err != nil {
		logFatal(err)
	}
	return item_state.length
}

func (self *State) getItemById(target_id uint32, type_name string) (c Content) {
	var state_slice *List

	state_slice, err := self.getStateSliceByName(type_name)
	if err != nil {
		logFatal(err)
	}
	var current_node *ListNode = state_slice.root
	for current_node != nil {
		c = current_node.NodeContent
		if c.getId() == target_id {
			return
		}
		current_node = current_node.Next
	}
	return nil
}

func (self *State) getItemByToString(target string, type_name string) Content {
	var target_slice *List
	target_slice, err := self.getStateSliceByName(type_name)
	if err != nil {
		logFatal(err)
	}
	c := target_slice.exists(target)
	return c
}

func (self *State) getItemsByToStringPrefix(search_value string, type_name string) string {
	var target_slice *List
	target_slice, err := self.getStateSliceByName(type_name)
	if err != nil {
		logFatal(err)
	}
	var results string = "["
	for current_node := target_slice.root; current_node != nil; current_node = current_node.Next {
		if strings.HasPrefix(current_node.NodeContent.toString(), search_value) {
			results += current_node.NodeContent.toPartialJson() + ","
		}
	}
	if len(results) > 1 {
		return results[:len(results)-1] + "]"
	} else {
		return "[]"
	}
}

func (self *State) getItemsAsJson(type_name string) string {
	// only gets id, username and name
	var items_json string = "["
	var type_state *List
	type_state, err := self.getStateSliceByName(type_name)
	if err != nil {
		logFatal(err)
	}
	var current_item Content
	current_node := type_state.root
	for current_node != nil {
		current_item = current_node.NodeContent
		items_json += current_item.toPartialJson()
		if current_node.Next != nil {
			items_json += ","
		}
		current_node = current_node.Next
	}
	return items_json + "]"
}

func (self *State) loadState() error {
	var err error
	if pathExists(SERVER_DATA) {
		if err = self.loadStateSlice(USERS_STATE); err != nil {
			return err
		}
		if err = self.loadStateSlice(CLIENTS_STATE); err != nil {
			return err
		}
		if err = self.loadStateSlice(VEHICLES_STATE); err != nil {
			return err
		}
		return err
	} else {
		return fmt.Errorf("'%s' no such file or directory", SERVER_DATA)
	}
}

func (self *State) loadStateSlice(type_name string) error {
	var load_path string = self.composeSaveFile(type_name)
	if pathExists(load_path) {
		type_state, err := self.getStateSliceByName(type_name)
		if err != nil {
			logFatal(err)
		}
		f, err := ioutil.ReadFile(load_path)
		if err != nil {
			return err
		}
		var new_item Content
		for _, u := range strings.Split(string(f), "*") {
			new_item = self.getTypeByName(type_name)
			if new_item.load(u) == nil {
				type_state.append(new_item)
			}
		}
		return err
	} else {
		fmt.Printf("No such file or directory: %s\n", load_path)
		return nil
	}
}

func (self *State) saveStateToLocalStorage() error {
	var err error
	fmt.Println("Saving server state...")
	if pathExists(SERVER_DATA) {
		if err = self.saveStateSlice(USERS_STATE); err != nil {
			return err
		}
		if err = self.saveStateSlice(CLIENTS_STATE); err != nil {
			return err
		}

	} else {
		err = fmt.Errorf("'%s' no such file or directory", SERVER_DATA)
	}
	return err
}

func (self *State) saveStateSlice(type_name string) error {
	type_state, err := self.getStateSliceByName(type_name)
	if err != nil {
		logFatal(err)
	}
	var rstring string
	var c Content
	for node := type_state.root; node != nil; node = node.Next {
		c = node.NodeContent
		rstring += fmt.Sprintf("%s*", c.toRstring())
	}
	fmt.Printf("%s saved!\n", type_name)
	return ioutil.WriteFile(self.composeSaveFile(type_name), []byte(rstring), 0755)
}

func (self *State) usernameExists(username string) bool {
	c := self.users.exists(username)
	return c != nil
}
