package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strconv"
)

type Resp struct {
	RC   int         `json:"rc"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func (data *Resp) Encode() []byte {

	buf, _ := json.Marshal(data)

	return buf
}

type UserDefPanic struct {
	Code int
	Msg  string
}

type APIHandler func(http.ResponseWriter, *http.Request) interface{}

func entryHandler(f APIHandler) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		defer panicHandle(w, req)

		action := req.Header.Get("Action")
		token := req.Header.Get("Token")
		userid := req.FormValue("userid")
		uid, _ := strconv.Atoi(userid)
		if action != "register" && action != "login" {
			//check token

			if !GobalHub.isValidToken(uid, token) {
				p := UserDefPanic{Code: 0, Msg: "token error"}
				panic(p)
			}
		}

		ret := f(w, req)

		var resp = Resp{0, "success", ret}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp.Encode())
	}
}

func panicHandle(w http.ResponseWriter, r *http.Request) {

	if err := recover(); err != nil {
		if user, ok := err.(UserDefPanic); ok {

			var resp = Resp{1, user.Msg, nil}
			w.Header().Set("Content-Type", "application/json")
			if user.Code == 0 {

				w.WriteHeader(http.StatusBadRequest)
			} else {

				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(resp.Encode())
		} else {

			LogError("handler panic %v", err)
			LogError("crash stack: %v", string(debug.Stack()))
		}

	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	LogDebug(r.URL.String())
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}
