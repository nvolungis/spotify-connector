package main

import (
	"fmt"
	"github.com/badfortrains/spotcontrol"
  "github.com/gorilla/mux"
  "github.com/rs/cors"
	"strings"
  "net/http"
  "time"
  "flag"
)

const defaultdevicename = "SpotControl"
var username string
var password string
var appkey string


func login() *spotcontrol.SpircController{

  devicename := defaultdevicename
	return spotcontrol.Login(username, password, appkey, devicename)
}


func chooseFirstDevice(controller *spotcontrol.SpircController) string{
	devices := controller.ListDevices()
  for len(devices) == 0 {
    time.Sleep(100 * time.Millisecond)
	  devices = controller.ListDevices()
  }
  return devices[0].Ident
}


func handlePlay(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("Sending Play")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  sController := login()
  ident := chooseFirstDevice(sController)
  fmt.Fprintf(w, "Sending Play", ident)
	sController.SendPlay(ident)
}


func handlePause(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  sController := login()
  ident := chooseFirstDevice(sController)
  fmt.Fprintf(w, "Sending Pause", ident)
	sController.SendPause(ident)
}


func handleLoad(w http.ResponseWriter, r *http.Request) {
  if (r.Method == "OPTIONS") {
    w.Header().Set("Access-Control-Allow-Origin", "*")
  }else {
    sController := login()
    ident := chooseFirstDevice(sController)
    vars := mux.Vars(r)
    songIds := strings.Split(vars["songId"], ",")
    fmt.Fprintf(w, "Sending load", songIds, ident)
    sController.LoadTrack(ident, songIds)
  }
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
}


func main() {
  _username := flag.String("username", "", "spotify username")
  _password := flag.String("password", "", "spotify password")
  _appkey := flag.String("appkey", "./spotify_appkey.key", "spotify appkey file path")
  flag.Parse()

  username = *_username
  password = *_password
  appkey = *_appkey

  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/load/{songId}", handleLoad)
  router.HandleFunc("/play", handlePlay)
  router.HandleFunc("/pause", handlePause)
  router.HandleFunc("",handleOptions).Methods("OPTIONS")

  handler := cors.Default().Handler(router)
  http.ListenAndServe(":8080", handler)
}
