package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func newRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/time", getTimeHandler).Methods("GET")
	router.HandleFunc("/ws", getWebSocketHandler)

	// /assets/ Page
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	return router
}

func main() {
	router := newRouter()

	http.ListenAndServe(":8080", router)
}

func getTimeHandler(w http.ResponseWriter, r *http.Request) {

	timeBytes, err := json.Marshal(getTimeNow())

	// If there is an error, print it to the console, and return a server
	// error response to the user
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If all goes well, write the JSON to the response
	w.Write(timeBytes)
}

func getWebSocketHandler(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Origin") != "http://"+r.Host {
		//"ws://" + location.host + "/ws"
		http.Error(w, "Origin not allowed", 403)
		return
	}

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go timeRoutine(conn)

}

func timeRoutine(conn *websocket.Conn) {

	t := time.NewTicker(time.Second * 3)
	for {

		timeString := getTimeNow()

		if err := conn.WriteJSON(timeString); err != nil {
			fmt.Println(err)
		}

		<-t.C
	}
}

func getTimeNow() string {
	timeNow := time.Now()
	formattedTime := timeNow.Format("Mon Jan 02 15:04:05 MST 2006")

	return formattedTime
}
