package main

import _ "github.com/astaxie/beego/session/memcache"
import (
	"net/http"
	"github.com/astaxie/beego"
	"html/template"
)
func count(w http.ResponseWriter, r *http.Request){
	sess,_ := beego.GlobalSessions.SessionStart(w,r)
	ct := sess.Get("countnum")
	if ct == nil{
		sess.Set("countnum",1)
	}else{
		sess.Set("countnum",(ct.(int) + 1))
	}
	t, _ := template.ParseFiles("count.gtpl")
	w.Header().Set("Content-Type","text/html")
	t.Execute(w, sess.Get("countnum"))
}
func main() {

	http.HandleFunc("/",count)
	http.ListenAndServe("localhost","7071")
}
