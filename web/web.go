package web

import (
	//Third party
	"github.com/fatih/color"
	//
	"net/http"
	"io/ioutil"
	"time"
	"crypto/subtle"
	"config"
	"context"
	"strconv"
)

var logdir string
//const auth401 = http.Error("Unauthorized",401)
type WebServer struct {
	server *http.Server
	config config.WebConfig
}

func  NewWebServer (config config.WebConfig) *WebServer{
	defaultServer := &http.Server{
		ReadTimeout: 10*time.Second,
		Addr: config.ListenIP+":"+strconv.Itoa(config.ListenWebPort),
		//Handler: http.HandlerFunc(WebHandler),
	}
	ws := &WebServer {
		server: defaultServer,
		config: config,
	}
	return ws
}

func WebHandler(w http.ResponseWriter, r *http.Request){
	data, _ := ioutil.ReadFile("./logs/log.txt")
	w.Write(data)
}

func BasicAuth (handler http.HandlerFunc, username string, password string) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request){
		//before := time.Now()
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user),[]byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass),[]byte(password)) != 1 {
		//secondspassed := time.Now().Sub(before)
		w.Header().Set("WWW-Authenticate","Basic realm=SOSAI")
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized.\n"))
		return
		}
	handler(w,r)
	}
}
	
func (ws *WebServer) Run() {
	color.Cyan("[!] Starting Web server at http://%s",ws.config.ListenIP)
	http.HandleFunc(ws.config.LogDirectory,BasicAuth(WebHandler,ws.config.Username,ws.config.Password))
	ws.server.ListenAndServe()
}
func (ws *WebServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	return ws.server.Shutdown(ctx)
}
