package main

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type AppContext struct {
	c valkey.Client
}

type Status int

const (
	Done = iota
	Started
	New
)

type Todo struct {
	Description string
	Status      Status
}

var ctx = context.Background()

func (appCtx AppContext) listKeys(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "You shall not pass! Only GET supported here!!! :o", http.StatusMethodNotAllowed)
		return
	}

	// writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")

	cmd := appCtx.c.B().Keys().Pattern("*").Build()
	keys := appCtx.c.Do(request.Context(), cmd)

	karr, err := keys.AsStrSlice()
	if err != nil {
		http.Error(writer, "{ \"status\": \"listKeys returned nothing\" }", http.StatusMethodNotAllowed)
		return
	}

	keysInJson, _ := json.Marshal(karr)

	writer.Write([]byte(keysInJson))
}

func (appCtx AppContext) addItem(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "You shall not pass! Only POST supported here!!! :o", http.StatusMethodNotAllowed)
		return
	}
	// Parse the request to create a Todo entry
	var todo Todo

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check that the body can be unmarshalled correctly
	err = json.Unmarshal(body, &todo)
	if err != nil {

		http.Error(writer, "Content cannot be unmarshalled, check your json element", http.StatusBadRequest)
		return
	}

	defer request.Body.Close()

	// Quietly produce an uuid for every message going into the valkey
	uuid := uuid.New().String()

	bodyString := string(body)

	if bodyString == "" {
		http.Error(writer, "Body is empty, cannot proceed", http.StatusBadRequest)
		return
	}

	// Store the struct as an item
	cmd := appCtx.c.B().Set().Key(uuid).Value(bodyString).Build()
	appCtx.c.Do(request.Context(), cmd)

	fmt.Println("New item UUID added: ", uuid)
	fmt.Println("Body of new item: ", bodyString)

	writer.WriteHeader(http.StatusCreated)
}

func (appCtx AppContext) getItem(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "You shall not pass! Only GET supported here!!! :o", http.StatusMethodNotAllowed)
		return
	}

	uuid := request.URL.Query().Get("uuid")
	fmt.Println("Retrieving item with UUID: ", uuid)

	// var todo Todo
	cmd := appCtx.c.B().Get().Key(uuid).Build()
	entryInString, _ := appCtx.c.Do(request.Context(), cmd).ToString()

	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(entryInString))
}

func (cx AppContext) emptyPage(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Empty page mother fucker", http.StatusMethodNotAllowed)
}

func getEnvironment() (string, string) {
	server := cmp.Or(os.Getenv("VALKEY_SERVER"), "localhost")
	password := cmp.Or(os.Getenv("VALKEY_PASSWORD"), "TestingThisShitYo")
	if password == "" {
		fmt.Println("Environment variable VALKEY_PASSWORD is not set")
		os.Exit(1)
	}
	return server, password
}

func setupValkey() valkey.Client {
	valkey_server, valkey_password := getEnvironment()
	fmt.Println("Valkey server:", valkey_server)

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkey_server + ":6379"},
		Password:    valkey_password,
	})

	if err != nil {
		fmt.Println("Failed to create Valkey client")
		os.Exit(1)
	}
	fmt.Println("Valkey client created successfully")
	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		fmt.Println("Failed to ping Valkey server:", err)
		os.Exit(1)
	}
	fmt.Println("Valkey server is reachable")

	return client
}

func getInterNetworkInterface() string {
	var ifaces []net.Interface
	var err error
	if ifaces, err = net.Interfaces(); err != nil {
		fmt.Println("Assigning interfaces from net.Interfaces() failed")
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Couldn't obtain addresses from the interface")
		}
		if iface.Name == "lo" {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				ip := ipnet.IP
				if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() {
					return ip.String()
				}
			}
		}
	}
	return ""
}

func main() {
	client := setupValkey()
	appCtx := &AppContext{
		c: client,
	}

	http.HandleFunc("/", appCtx.emptyPage)
	http.HandleFunc("/listKeys", appCtx.listKeys)
	http.HandleFunc("/getItem", appCtx.getItem)
	http.HandleFunc("/addItem", appCtx.addItem)

	iface := getInterNetworkInterface()
	fmt.Println("Trying to listen to: " + iface)
	http.ListenAndServe(iface+":3000", nil)
}
