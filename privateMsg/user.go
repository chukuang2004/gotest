package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type User struct {
	hub *Hub
}

func (self *User) Register(w http.ResponseWriter, r *http.Request) interface{} {

	id := r.FormValue("id")
	pwd := r.FormValue("pwd")

	if id == "" || pwd == "" {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	LogDebug("Register, args: %s,%s", id, pwd)

	var sql string = "INSERT IGNORE INTO `Users`(`account`,`pwd`) " +
		"VALUES(?, ?)"

	stmt, err := DB.Prepare(sql)
	if err != nil {

		LogError("prepare db fail, sql#%s,error:%s", sql, err.Error())
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}
	defer stmt.Close()

	ret, err := stmt.Exec(id, pwd)
	if err != nil {
		LogError("store db fail : %s", err.Error())
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	cnt, _ := ret.RowsAffected()

	if cnt == 0 {
		p := UserDefPanic{Code: 0, Msg: "id repeat"}
		panic(p)
	}

	return nil
}

// login(email,nickname,pwd)
func (self *User) Login(w http.ResponseWriter, r *http.Request) interface{} {

	id := r.FormValue("id")
	pwd := r.FormValue("pwd")

	if id == "" || pwd == "" {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	LogDebug("Login, args: %s,%s", id, pwd)

	sql := fmt.Sprintf("SELECT userid FROM `Users` WHERE account = \"%s\" AND pwd = %s", id, pwd)
	rows, err := DB.Query(sql)
	if err != nil {
		LogError("read db fail, sql#%s,error:%s", sql, err.Error())
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	defer rows.Close()

	var userid int
	if rows.Next() {

		err = rows.Scan(&userid)
		if err != nil {
			LogError("read db fail, %s", err.Error())
			p := UserDefPanic{Code: 1, Msg: "internal error"}
			panic(p)
		}
	}

	buf := md5.Sum(bytes.NewBufferString(time.Now().String()).Bytes())
	token := base64.StdEncoding.EncodeToString(buf[:])
	self.hub.BindToken(userid, token)

	type Result struct {
		Userid int    `json:"userid"`
		Token  string `json:"token"`
	}

	return Result{Userid: userid, Token: token}
}

// addFriend(userid,userid)//no permit
func (self *User) AddFriend(w http.ResponseWriter, r *http.Request) interface{} {

	userid := r.FormValue("userid")
	peerid := r.FormValue("peerid")

	if userid == "" || peerid == "" {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	flag := isUserExist(userid)
	if flag == false {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}
	flag = isUserExist(peerid)
	if flag == false {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}
	LogDebug("AddFriend, args: %s %s", userid, peerid)

	var err error = nil
	for i := 0; i < 3; i++ {
		err = changeFriends(userid, peerid, true)
		if err == nil {
			break
		} else {
			LogError("%s", err.Error())
		}
	}
	if err == nil {
		for i := 0; i < 3; i++ {
			err = changeFriends(peerid, userid, true)
			if err == nil {
				break
			} else {
				LogError("%s", err.Error())
			}
		}
	} else {
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	if err != nil {
		//roll back
		for i := 0; i < 3; i++ {
			err = changeFriends(userid, peerid, false)
			if err == nil {
				break
			} else {
				LogError("%s", err.Error())
			}
		}
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	return nil
}

// delfriend(userid,userid)
func (self *User) Delfriend(w http.ResponseWriter, r *http.Request) interface{} {

	userid := r.FormValue("userid")
	peerid := r.FormValue("peerid")

	if userid == "" || peerid == "" {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	flag := isUserExist(userid)
	if flag == false {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}
	flag = isUserExist(peerid)
	if flag == false {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}
	LogDebug("Delfriend, args: %s %s", userid, peerid)

	var err error = nil
	for i := 0; i < 3; i++ {
		err = changeFriends(userid, peerid, false)
		if err == nil {
			break
		} else {
			LogError("%s", err.Error())
		}
	}
	if err == nil {
		for i := 0; i < 3; i++ {
			err = changeFriends(peerid, userid, false)
			if err == nil {
				break
			} else {
				LogError("%s", err.Error())
			}
		}
	} else {
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	if err != nil {
		//roll back
		for i := 0; i < 3; i++ {
			err = changeFriends(userid, peerid, true)
			if err == nil {
				break
			} else {
				LogError("%s", err.Error())
			}
		}
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	return nil
}

// friendsList(userid)(userid,nickname,newMsgCnt)
func (self *User) FriendsList(w http.ResponseWriter, r *http.Request) interface{} {

	userid := r.FormValue("userid")

	if userid == "" {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	flag := isUserExist(userid)
	if flag == false {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	LogDebug("FriendsList, args: %s", userid)

	sql := fmt.Sprintf("SELECT friends FROM `Users` WHERE userid = %s", userid)
	rows, err := DB.Query(sql)
	if err != nil {
		LogError("read db fail, sql#%s,error:%s", sql, err.Error())
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	defer rows.Close()

	type Friend struct {
		Userid  int    `json:"userid"`
		Account string `json:"account"`
	}

	result := make([]Friend, 0)
	var friends string
	if rows.Next() {

		err = rows.Scan(&friends)
		if err != nil {
			LogError("read db fail, %s", err.Error())
			p := UserDefPanic{Code: 1, Msg: "internal error"}
			panic(p)
		}

		f := strings.Trim(friends, ",")
		if f == "" {
			return result
		}

		sql := fmt.Sprintf("SELECT userid,account from Users where userid in (%s)", f)
		rows1, err := DB.Query(sql)
		if err != nil {
			LogError("read db fail, sql#%s,error:%s", sql, err.Error())
			p := UserDefPanic{Code: 1, Msg: "internal error"}
			panic(p)
		}

		defer rows1.Close()
		var uid int
		var account string
		for rows1.Next() {

			err = rows1.Scan(&uid, &account)
			if err != nil {
				LogError("read db fail, %s", err.Error())
				p := UserDefPanic{Code: 1, Msg: "internal error"}
				panic(p)
			}

			result = append(result, Friend{Userid: uid, Account: account})
		}
	}

	return result
}

// getHistoryMsg(time,cnt)
func (self *User) GetHistoryMsg(w http.ResponseWriter, r *http.Request) interface{} {

	userid := r.FormValue("userid")

	if userid == "" {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	flag := isUserExist(userid)
	if flag == false {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}
	LogDebug("GetHistoryMsg, args: %s", userid)

	return nil
}

func (self *User) DelMsg(w http.ResponseWriter, r *http.Request) interface{} {

	userid := r.FormValue("userid")

	if userid == "" {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}

	flag := isUserExist(userid)
	if flag == false {
		p := UserDefPanic{Code: 0, Msg: "wrong args"}
		panic(p)
	}
	LogDebug("DelMsg, args: %s", userid)

	return nil
}

func changeFriends(userid, peerid string, isAdd bool) error {

	sql := fmt.Sprintf("SELECT userid,friends,version FROM `Users` WHERE userid = %s", userid)
	rows, err := DB.Query(sql)
	if err != nil {
		return fmt.Errorf("read db fail, sql#%s,error:%s", sql, err.Error())
	}

	defer rows.Close()

	var (
		uid     int
		friends string
		version int
	)

	if rows.Next() {

		err = rows.Scan(&uid, &friends, &version)
		if err != nil {
			return fmt.Errorf("read db fail, %s", err.Error())
		}
	}

	if isAdd {
		if friends != "" {
			if !strings.Contains(friends, ","+peerid+",") {
				friends = friends + peerid + ","
			} else {
				return nil
			}
		} else {
			friends = fmt.Sprintf(",%s,", peerid)
		}
	} else {
		if friends != "" {
			if strings.Contains(friends, ","+peerid+",") {

				friends = strings.Replace(friends, ","+peerid+",", ",", -1)
				if friends == "," {
					friends = ""
				}
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	sql = "UPDATE `Users` SET `friends` = ?, `version`=? where `userid`=? AND `version`=?"

	stmt, err := DB.Prepare(sql)
	if err != nil {

		return fmt.Errorf("prepare db fail, sql#%s,error:%s", sql, err.Error())
	}
	defer stmt.Close()

	ret, err := stmt.Exec(friends, version+1, uid, version)
	if err != nil {
		return fmt.Errorf("store db fail : %s", err.Error())
	}
	cnt, _ := ret.RowsAffected()
	if cnt != 1 {
		return fmt.Errorf("update conflict")
	}

	return nil
}

func isUserExist(userid string) bool {

	sql := fmt.Sprintf("SELECT userid FROM `Users` WHERE userid = %s", userid)
	rows, err := DB.Query(sql)
	if err != nil {
		LogError("read db fail, sql#%s,error:%s", sql, err.Error())
		p := UserDefPanic{Code: 1, Msg: "internal error"}
		panic(p)
	}

	defer rows.Close()

	var uid int = 0
	if rows.Next() {

		err = rows.Scan(&uid)
		if err != nil {
			LogError("read db fail, %s", err.Error())
			p := UserDefPanic{Code: 1, Msg: "internal error"}
			panic(p)
		}
	}

	if uid != 0 {
		return true
	} else {
		return false
	}
}
