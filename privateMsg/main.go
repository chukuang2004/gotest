package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var configPath = flag.String("configPath", "config.json", "config file path")

// register
// friendsList(userid)(userid,nickname,newMsgCnt)
// addFriend(userid,userid)//no permit
// delfriend(userid,userid)
// getHistoryMsg(time,cnt)
// delMsg()
// login(email,nickname,pwd)

// getMsg/sendMsg()//addfrind when sendMsg

var GobalHub *Hub = nil

func main() {
	flag.Parse()

	SetFlags(log.Ldate|log.Lmicroseconds|log.Lshortfile, "[Webim] ")

	_, err := InitConfiguration(*configPath)
	if err != nil {
		LogError("config error %s\n", err.Error())
		return
	}

	SetLogLevelString(GetLogLevel())

	err = InitDB()
	if err != nil {
		LogError("init db error %s", err.Error())
		return
	}

	GobalHub = InitHub()

	user := User{hub: GobalHub}
	r := mux.NewRouter()
	r.HandleFunc("/user", entryHandler(user.Register)).Methods("POST").Headers("Action", "register", "Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/user", entryHandler(user.Login)).Methods("POST").Headers("Action", "login", "Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/user", entryHandler(user.FriendsList)).Methods("POST").Headers("Action", "friendsList", "Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/user", entryHandler(user.AddFriend)).Methods("POST").Headers("Action", "addFriend", "Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/user", entryHandler(user.Delfriend)).Methods("POST").Headers("Action", "delFriend", "Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/user", entryHandler(user.GetHistoryMsg)).Methods("POST").Headers("Action", "getHistoryMsg", "Content-Type", "application/x-www-form-urlencoded")
	r.HandleFunc("/user", entryHandler(user.DelMsg)).Methods("POST").Headers("Action", "delMsg", "Content-Type", "application/x-www-form-urlencoded")

	r.HandleFunc("/", serveHome)
	r.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {

		ServeWS(GobalHub, w, r)

	})

	err = http.ListenAndServe(GetListenAddr(), r)
	if err != nil {
		LogFatal("ListenAndServe: ", err)
	}
}
