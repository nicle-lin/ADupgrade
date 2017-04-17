package main

import "net/http"

func response(w http.ResponseWriter, req *http.Request){
	w.Write([]byte("hello world and how old are you!"))
	req.Body
}
func main() {
	http.HandleFunc("/index.html",response)
	http.HandleFunc("/",response)
	http.ListenAndServe(":5001",nil)
}
