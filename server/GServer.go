package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const TYPESDATA = "./operational_data/types_data.json"

type Server struct {
	port         int
	state        *State
	sessions     map[uint64]*User
	types_data   map[string]interface{}
	ok           []byte
	bad_response []byte
}

func (self *Server) init(port int) {
	self.port = port
	self.state = new(State)
	self.state.init()
	self.sessions = make(map[uint64]*User)
	self.ok = []byte("ok")
	self.bad_response = []byte("error")
}

func (self *Server) composeResponse(response_value string) []byte {
	return []byte(fmt.Sprintf("{\"response\": %s}", response_value))
}

func (self *Server) composeErrorCode(error_code int) []byte {
	return []byte(fmt.Sprintf("{\"code\": %d}", error_code))
}

func (self *Server) createUser(r *http.Request) *User {
	username := r.FormValue("username")
	name := r.FormValue("name")
	middlename := r.FormValue("middlename")
	lastname := r.FormValue("lastname")
	address := r.FormValue("address")
	password := r.FormValue("password")
	phone := r.FormValue("phone")
	group := r.FormValue("user-group")
	if middlename != "" || lastname != "" {
		name = composeName(name, middlename, lastname)
	}
	return createUser(self.state.getNewItemId(USERS_STATE), username, name, password, address, phone, group)
}

func (self *Server) createClient(request *http.Request) *Client {
	var new_client *Client = new(Client)
	new_client.id = uint32(self.state.getNewItemId(CLIENTS_STATE))
	new_client.name = request.FormValue("name")
	new_client.address = request.FormValue("address")
	new_client.phone = request.FormValue("phone")
	return new_client
}

func (self *Server) createVehicle(request *http.Request) *Vehicle {
	var new_vehicle *Vehicle = new(Vehicle)
	new_vehicle.id = uint32(self.state.getNewItemId(VEHICLES_STATE))
	new_vehicle.plates = request.FormValue("plates")
	new_vehicle.model = request.FormValue("model")
	new_vehicle.brand = request.FormValue("brand")
	new_vehicle.release = request.FormValue("release")
	new_vehicle.client = uint32(stringToInt(request.FormValue("client")))
	return new_vehicle
}

func (self *Server) createResponse(response string) []byte {
	return []byte(fmt.Sprintf("{\"response\":\"%s\"}", response))
}

func (self *Server) createSession(user *User, remote_address string) uint64 {
	var session_uuid uint64 = shaAsInt64(remote_address + user.username)
	self.sessions[session_uuid] = user
	fmt.Printf("new-session-uuid: %d\n", session_uuid)
	return session_uuid
}

func (self *Server) enableCors(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Headers", "X-sk")
		response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
		if request.Method == http.MethodOptions {
			response.WriteHeader(200)
			response.Write(self.ok)
			return
		}

		handler(response, request)
	}
}

func (self *Server) getPortStr() string {
	return fmt.Sprintf(":%d", self.port)
}

func (self *Server) getRequestUser(request *http.Request) (*User, error) {
	var string_key string = request.Header.Get("X-sk")
	if string_key != "" {
		var sk uint64 = stringToUint64(string_key)
		if u, ok := self.sessions[sk]; ok {
			return u, nil
		} else {
			fmt.Printf("Invalid session key: %d\n", sk)
			return nil, fmt.Errorf("Invalid session key: %d\n", sk)
		}
	} else {
		fmt.Println("Recived header with missing session_key")
		return nil, fmt.Errorf("missing header")
	}
}

func (self *Server) logUser(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		var username string = request.URL.Query().Get("username")
		if u := self.state.getItemByToString(username, USERS_STATE).(*User); u != nil {
			var raw_password string = request.URL.Query().Get("password")
			var password uint64 = shaAsInt64(raw_password)
			if password == u.password {
				fmt.Printf("Session started for user %s with remote_address:%s\n", u.username, request.RemoteAddr)
				session_uuid := self.createSession(u, request.RemoteAddr)
				response.WriteHeader(http.StatusOK)
				response.Header().Set("Content-Type", "application/json")
				_, _ = response.Write(self.composeResponse(fmt.Sprintf("\"%d\"", session_uuid)))
			} else {
				fmt.Printf("Wrong password for user: %s\n", u.username)
				response.WriteHeader(http.StatusNonAuthoritativeInfo)
				response.Header().Set("Content-Type", "application/json")
				_, _ = response.Write(self.composeErrorCode(1))
			}
		} else {
			// user doesnt exists
			fmt.Printf("user doesnt exists: %s\n", username)
			response.WriteHeader(http.StatusBadRequest)
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write(self.composeErrorCode(0))
		}
	} else {
		self.setNotImplemented(response, request)
	}
}

func (self *Server) loadTypesData() error {
	var types_bytes []byte
	var types_data map[string]interface{} = make(map[string]interface{})
	types_bytes, err := ioutil.ReadFile(TYPESDATA)
	if err != nil {
		fmt.Printf("Couldnt load types_data because the file '%s' is not valid or doesnt exists\n", TYPESDATA)
		return err
	}
	err = json.Unmarshal(types_bytes, &types_data)
	if err != nil {
		return err
	}
	self.types_data = types_data
	return err
}

func (self *Server) handleClients(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		var target string = request.URL.Query().Get("id")
		fmt.Printf("GET request for client %s...", target)
		if target == "*" {
			// request for all clients
			response.WriteHeader(200)
			response.Write([]byte(self.state.getItemsAsJson(CLIENTS_STATE)))
		} else {
			// request for a single client, will return all its data
			var string_id string = request.URL.Query().Get("id")
			if string_id != "" {
				var id uint32 = uint32(stringToInt(string_id))
				var client Content = self.state.getItemById(id, CLIENTS_STATE)
				if client != nil {
					fmt.Println("Succes!")
					response.WriteHeader(200)
					response.Write([]byte(client.toJson()))
				} else {
					// client doesnt exists
					fmt.Println("Failed: user doesnt exists")
					response.WriteHeader(404)
					response.Write(self.bad_response)
				}
			} else {
				response.WriteHeader(406)
				response.Write(self.bad_response)
			}
		}
	case http.MethodPost:
		// Registering Client
		var new_client *Client = self.createClient(request)
		fmt.Printf("POST request for client %s...", new_client.name)
		if new_client.name != "" {
			self.state.addItemToState(new_client, CLIENTS_STATE)
			if err := self.state.saveStateSlice(CLIENTS_STATE); err != nil {
				logFatal(err)
			}
			fmt.Println("Success!")
			response.WriteHeader(200)
			response.Write(self.ok)
		} else {
			fmt.Println("Failed: name for client cannot be empty")
			response.WriteHeader(406)
			response.Write(self.bad_response)
		}
	case http.MethodPatch:
		// UPDATE: restricted
		fmt.Print("UPDATE request for client with id ")
		if self.isUserAdmin(request) {
			var string_id string = request.FormValue("id")
			fmt.Printf("'%s'...", string_id)
			if string_id != "" {
				var target_client *Client = self.state.getItemById(uint32(stringToInt(string_id)), CLIENTS_STATE).(*Client)
				if target_client == nil {
					fmt.Println("Failed: client wasnt found")
					response.WriteHeader(404)
					response.Write(self.bad_response)
					return
				}
				self.updateClient(target_client, request)
				self.state.saveStateSlice(CLIENTS_STATE)
				fmt.Println("Succes!")
				response.WriteHeader(200)
				response.Write(self.ok)
			} else {
				fmt.Println("Failed: key is not vaild")
				response.WriteHeader(406)
				response.Write(self.bad_response)
			}
		} else {
			fmt.Println("Failed: User is not authorized")
			response.WriteHeader(401)
			response.Write(self.bad_response)
		}
	case http.MethodDelete:
		// DELETE: restricted
		fmt.Print("DELETE request for client ")
		if self.isUserAdmin(request) {
			var string_id string = request.FormValue("id")
			fmt.Printf("with id '%s'...", string_id)
			if string_id != "" {
				var id uint32 = uint32(stringToInt(string_id))
				err := self.state.deleteItemById(id, CLIENTS_STATE)
				if err != nil {
					fmt.Println("Failed: client wasnt found")
					response.WriteHeader(404)
					response.Write(self.bad_response)
					return
				}
				self.state.saveStateSlice(CLIENTS_STATE)
				fmt.Print("Success!")
				response.WriteHeader(200)
				response.Write(self.ok)
			}
		} else {
			fmt.Println("Failed: user is not authorized.")
			response.WriteHeader(401)
			response.Write(self.bad_response)
		}
	case http.MethodOptions:
		response.WriteHeader(200)
		_, _ = response.Write(self.ok)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = response.Write(self.bad_response)
	}
}

func (self *Server) handleVehicles(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		var target string = request.URL.Query().Get("id")
		fmt.Printf("GET request for vehicle %s...", target)
		if target == "*" {
			// request for all clients
			response.WriteHeader(200)
			response.Write([]byte(self.state.getItemsAsJson(VEHICLES_STATE)))
		} else {
			// request for a single client, will return all its data
			if target != "" {
				var id uint32 = uint32(stringToInt(target))
				var vehicle Content = self.state.getItemById(id, VEHICLES_STATE)
				if vehicle != nil {
					fmt.Println("Succes!")
					response.WriteHeader(200)
					response.Write([]byte(vehicle.toJson()))
				} else {
					// client doesnt exists
					fmt.Println("Failed: user doesnt exists")
					response.WriteHeader(404)
					response.Write(self.bad_response)
				}
			} else {
				response.WriteHeader(406)
				response.Write(self.bad_response)
			}
		}
	case http.MethodPost:
		// Registering Client
		var new_vehicle *Vehicle = self.createVehicle(request)
		fmt.Printf("POST request for vehicle %s...", new_vehicle.plates)
		if new_vehicle.plates != "" {
			self.state.addItemToState(new_vehicle, VEHICLES_STATE)
			if err := self.state.saveStateSlice(VEHICLES_STATE); err != nil {
				logFatal(err)
			}
			fmt.Println("Success!")
			response.WriteHeader(200)
			response.Write(self.ok)
		} else {
			fmt.Println("Failed: plates for vehicle cannot be empty")
			response.WriteHeader(406)
			response.Write(self.bad_response)
		}
	case http.MethodPatch:
		// UPDATE: restricted
		fmt.Print("UPDATE request for vehicle with id ")
		var string_id string = request.FormValue("id")
		fmt.Printf("'%s'...", string_id)
		if string_id != "" {
			var target_vehicle *Vehicle = self.state.getItemById(uint32(stringToInt(string_id)), VEHICLES_STATE).(*Vehicle)
			if target_vehicle == nil {
				fmt.Println("Failed: vehicle wasnt found")
				response.WriteHeader(404)
				response.Write(self.bad_response)
				return
			}
			self.updateVehicle(target_vehicle, request)
			self.state.saveStateSlice(VEHICLES_STATE)
			fmt.Println("Succes!")
			response.WriteHeader(200)
			response.Write(self.ok)
		} else {
			fmt.Println("Failed: key is not vaild")
			response.WriteHeader(406)
			response.Write(self.bad_response)
		}
	case http.MethodDelete:
		// DELETE: restricted
		fmt.Print("DELETE request for vehicle ")
		if self.isUserAdmin(request) {
			var string_id string = request.FormValue("id")
			fmt.Printf("with id '%s'...", string_id)
			if string_id != "" {
				var id uint32 = uint32(stringToInt(string_id))
				err := self.state.deleteItemById(id, VEHICLES_STATE)
				if err != nil {
					fmt.Println("Failed: vehicle wasnt found")
					response.WriteHeader(404)
					response.Write(self.bad_response)
					return
				}
				self.state.saveStateSlice(VEHICLES_STATE)
				fmt.Print("Success!")
				response.WriteHeader(200)
				response.Write(self.ok)
			}
		} else {
			fmt.Println("Failed: user is not authorized.")
			response.WriteHeader(401)
			response.Write(self.bad_response)
		}
	case http.MethodOptions:
		response.WriteHeader(200)
		_, _ = response.Write(self.ok)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = response.Write(self.bad_response)
	}
}

func (self *Server) handleUsers(response http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case http.MethodGet:
		var target string = request.URL.Query().Get("id")
		if target == "*" {
			response.WriteHeader(200)
			_, _ = response.Write([]byte(self.state.getItemsAsJson(USERS_STATE)))
		} else {
			target_id := stringToInt(target)
			target_user := self.state.getItemById(uint32(target_id), USERS_STATE).(*User)
			if target_user != nil {
				response.WriteHeader(http.StatusOK)
				_, _ = response.Write([]byte(target_user.toJson()))
			} else {
				self.setBadResponse(response, request)
			}
		}
	case http.MethodPost:
		// Registering Users
		var requester *User
		requester, err := self.getRequestUser(request)
		if err != nil {
			self.setBadResponse(response, request)
			return
		}
		// if request.ParseForm() != nil {
		// 	os.Exit(0)
		// }
		// fmt.Printf("Data: %s\n", request.PostForm.Encode())
		fmt.Printf("POST request for %s by %s\n", request.FormValue("username"), requester.username)
		if !self.state.usernameExists(request.FormValue("username")) && requester.group == 0 {
			fmt.Printf("User doesnt exists: %v\n", !self.state.usernameExists(request.FormValue("username")))
			self.registerUser(response, request)
		} else if requester.group != 0 {
			fmt.Printf("User %s cannot register user with a user group=%d\n", requester.username, requester.group)
			response.WriteHeader(http.StatusForbidden)
			_, _ = response.Write(self.composeResponse("Permission denied"))
		} else {
			fmt.Printf("Username %s already exists\n", request.FormValue("username"))
			response.WriteHeader(http.StatusAlreadyReported)
			_, _ = response.Write(self.composeResponse("Username already in use"))
		}
	case http.MethodPatch:
		// update user data
		// verifying user permissions
		if self.isUserAdmin(request) {
			var id uint32 = uint32(stringToInt(request.FormValue("id")))
			var user *User = self.state.getItemById(id, USERS_STATE).(*User)
			fmt.Printf("PATCH request for user %s\n", user.username)
			self.updateUser(user, request)
			self.state.saveStateSlice(USERS_STATE)
			response.WriteHeader(200)
			_, _ = response.Write(self.ok)
		} else {
			fmt.Println("Permission denied")
			response.WriteHeader(http.StatusForbidden)
			_, _ = response.Write(self.composeResponse("Permission denied"))
		}
	case http.MethodDelete:
		if self.isUserAdmin(request) {
			var target_id uint32 = uint32(stringToInt(request.FormValue("id")))
			fmt.Printf("DELETE request for %d\n", target_id)
			if self.state.deleteItemById(target_id, USERS_STATE) == nil {
				self.state.saveStateSlice(USERS_STATE)
				response.WriteHeader(200)
				response.Write(self.ok)
			} else {
				response.WriteHeader(406) // Not acceptable
				response.Write(self.composeResponse("invalid id"))
			}
		}
	case http.MethodOptions:
		response.WriteHeader(200)
		_, _ = response.Write(self.ok)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = response.Write(self.bad_response)
	}
}

func (self *Server) handleUserPatch(request *http.Request) error {
	return nil
}

func (self *Server) handleType(response http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		var type_name string = request.URL.Query().Get("type_name")
		fmt.Println(request.URL.String())
		if type_data, ok := self.types_data[type_name]; ok {
			json_data, _ := json.Marshal(type_data)
			json_data = self.processTypeData(json_data)
			response.WriteHeader(http.StatusOK)
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write(json_data)
		} else {
			fmt.Printf("Type '%s' doesnt exists...\n", type_name)
			response.WriteHeader(http.StatusBadRequest)
			_, _ = response.Write(self.bad_response)
		}
	} else {
		response.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = response.Write(self.composeResponse("wrong_method"))
	}
}

func (self Server) isUserAdmin(request *http.Request) bool {
	var user *User
	user, err := self.getRequestUser(request)
	if err != nil {
		return false
	}
	return user.group == 0
}

func (self *Server) processTypeData(data []byte) []byte {
	var string_data string = string(data)
	regex, err := regexp.Compile("@[a-z]+@")
	if err != nil {
		logFatal(err)
	}
	var processable_string string = regex.FindString(string_data)
	if processable_string != "" {
		processable_string = strings.ReplaceAll(processable_string, "@", "")

		var options string = "[]"
		switch processable_string {
		case "clients":
			options = self.state.getItemsAsOptions(CLIENTS_STATE)
		}
		string_data = strings.Replace(string_data, fmt.Sprintf("\"@%s@\"", processable_string), options, 1)
	}
	fmt.Println(string_data)
	return []byte(string_data)

}

func (self *Server) setBadResponse(response http.ResponseWriter, r *http.Request) {
	response.WriteHeader(http.StatusBadRequest)
	_, _ = response.Write(self.createResponse("bad request"))
}

func (self *Server) setNotImplemented(response http.ResponseWriter, r *http.Request) {
	response.WriteHeader(http.StatusNotImplemented)
	_, _ = response.Write(self.bad_response)
}

func (self *Server) sessionExists(string_key string) bool {
	var sk uint64
	var err error
	sk, err = strconv.ParseUint(string_key, 0, 64)
	if err != nil {
		return false
	}
	_, exists := self.sessions[sk]
	return exists
}

func (self *Server) searchItem(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Serving search request")
	var search_type string = request.URL.Query().Get("type_name")
	var search_value string = request.URL.Query().Get("value")
	if search_type != "" && search_value != "" {
		var target string
		switch search_type {
		case "users":
			target = self.state.getItemsByToStringPrefix(search_value, USERS_STATE)
		case "clients":
			target = self.state.getItemsByToStringPrefix(search_value, CLIENTS_STATE)
		default:
			fmt.Println("Invalid type ", search_type)
			response.WriteHeader(501)
			response.Write(self.bad_response)
			return
		}
		if target == "[]" {
			fmt.Println("Not Found")
			response.WriteHeader(404)
		} else {
			fmt.Println("Success")
			response.WriteHeader(200)
		}
		response.Write([]byte(target))
	} else {
		fmt.Println("Not enought arguments for search request")
		response.WriteHeader(400)
		response.Write(self.bad_response)
	}
}

func (self *Server) registerUser(response http.ResponseWriter, request *http.Request) {
	var new_user *User = self.createUser(request)
	fmt.Printf("new user: %s\n", new_user)
	self.state.addItemToState(new_user, "users")
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write(self.ok)
}

func (self *Server) updateUser(user *User, request *http.Request) {
	var fake_user *User = self.createUser(request)
	if fake_user.username != user.username && fake_user.username != "" {
		user.username = fake_user.username
	}
	if fake_user.name != user.name && fake_user.name != "" {
		user.name = fake_user.name
	}
	if fake_user.password != 0 {
		user.password = fake_user.password
	}
	if fake_user.address != user.address && fake_user.address != "" {
		user.address = fake_user.address
	}
	if fake_user.phone != user.phone && fake_user.phone != "" {
		user.phone = fake_user.phone
	}
	if fake_user.group != user.group {
		user.group = fake_user.group
	}
	return
}

func (self *Server) updateClient(client *Client, request *http.Request) {
	var fake_client *Client = self.createClient(request)
	if fake_client.name != client.name && fake_client.name != "" {
		client.name = fake_client.name
	}
	if fake_client.address != client.address && fake_client.address != "" {
		client.address = fake_client.address
	}
	if fake_client.phone != client.phone && fake_client.phone != "" {
		client.phone = fake_client.phone
	}
}

func (self *Server) updateVehicle(vehicle *Vehicle, request *http.Request) {
	var fake_vehicle *Vehicle = self.createVehicle(request)
	if fake_vehicle.plates != vehicle.plates && fake_vehicle.plates != "" {
		vehicle.plates = fake_vehicle.plates
	}
	if fake_vehicle.model != vehicle.model && fake_vehicle.model != "" {
		vehicle.model = fake_vehicle.model
	}
	if fake_vehicle.brand != vehicle.brand && fake_vehicle.brand != "" {
		vehicle.brand = fake_vehicle.brand
	}
	if fake_vehicle.release != vehicle.release && fake_vehicle.release != "" {
		vehicle.release = fake_vehicle.release
	}
	if fake_vehicle.client >= 0 && fake_vehicle.client < uint32(self.state.clients.length) {
		vehicle.client = fake_vehicle.client
	}
}

func (self *Server) verifyAuthentication(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(response http.ResponseWriter, r *http.Request) {
		var session_key string = r.Header.Get("X-sk")
		if self.sessionExists(session_key) {
			fmt.Printf("Serving %s for %s\n", r.URL.String(), session_key)
			handler(response, r)
		} else {
			fmt.Printf("failed authentication_key: %s\n", session_key)
			response.WriteHeader(http.StatusNonAuthoritativeInfo)
			_, _ = response.Write(self.composeResponse("\"session key invalid\""))
		}
	}
}

func (self *Server) run() {
	if self.state.loadState() != nil || self.loadTypesData() != nil {
		fmt.Println("WARNING: server state or types couldnt be loaded couldnt be loaded!")
	}

	fmt.Printf("Users loaded: %d\n", self.state.users.length)
	fmt.Printf("Clients loaded: %d\n", self.state.clients.length)

	http.HandleFunc("/users", self.enableCors(self.verifyAuthentication(self.handleUsers)))
	http.HandleFunc("/clients", self.enableCors(self.verifyAuthentication(self.handleClients)))
	http.HandleFunc("/vehicles", self.enableCors(self.verifyAuthentication(self.handleVehicles)))
	http.HandleFunc("/login", self.enableCors(self.logUser))
	http.HandleFunc("/type", self.enableCors(self.verifyAuthentication(self.handleType)))
	http.HandleFunc("/search", self.enableCors(self.verifyAuthentication(self.searchItem)))

	fmt.Printf("Server running on 127.0.0.1%s\n", self.getPortStr())
	if http.ListenAndServe(self.getPortStr(), nil) != nil {
		panic(fmt.Errorf("Server error, sorry for the inconvinence"))
	}
}

func main() {
	var server *Server = new(Server)
	server.init(5000)
	server.run()
}
