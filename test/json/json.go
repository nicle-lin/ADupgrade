package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	ServerName string
	ServerIp string
}

type ServerSlice struct {
	Servers []Server
}
func main() {

	var Srv ServerSlice
	str := `{"Servers":[{"ServerName": "Shahai", "ServerIp":"127.0.0.1"},
	                  {"ServerName":"ShenZhen","ServerIp":"127.0.0.2"}]}`
	err := json.Unmarshal([]byte(str),&Srv)
	if err != nil{
		fmt.Println("err:",err)
		return
	}
	fmt.Println(Srv)
	fmt.Println(Srv.Servers[1].ServerIp)


	strByte := []byte(`{"Name":"wed","Age":6,"Parents":["father","mother"]}`)
	var f interface{}
	err = json.Unmarshal(strByte,&f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(f)

	fm ,ok:= f.(map[string]interface{})
	if ok {
		fmt.Println("good")
	}else {
		fmt.Println("bad")
	}
	fmt.Println("---------------")
	fmt.Println(fm["Name"])
	fmt.Println(fm["Parents"].([]interface{})[1])
	for k, v := range fm["Parents"].([]interface{}){
		fmt.Println(k,":",v)
	}

	Marshal()
}

func Marshal(){
	type Server struct {
		ServerName string `json:"server_name"`
		ServerIp string `json:omitempty`
	}

	type ServerSlice struct {
		Servers []Server `json:"servers"`
	}

	var S ServerSlice
	S.Servers = append(S.Servers,Server{ServerName:"ShaHai",ServerIp:"127.0.0.1"})
	S.Servers = append(S.Servers,Server{ServerName:"ShenZhen",ServerIp:"127.0.0.2"})
	S.Servers = append(S.Servers,Server{ServerName:"GuangZhou"})

	b, err := json.Marshal(S)
	if err != nil {
		fmt.Println("marshal error:",err)
		return
	}

	fmt.Println(string(b))
}
