package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)


type Step struct {
	StepDescription string
	StepOrderNum string
	StepName string
	StepParams interface{}
}


type Message struct {
	ThreadName string
	Resources []string
	CurrentStep Step
}




var _ws *websocket.Conn



//链接一次 记下来
func SetWs(ws *websocket.Conn){
	_ws=ws
}
func GetWs()*websocket.Conn{
	return _ws
}

func CloseWs(){
	if _ws!=nil{
		_ws.Close()
		fmt.Println("已经关闭当前ws")
	}
}

func SendMessage(m Message)error{
	data,err:=json.Marshal(&m)
	if err!=nil{
		panic(err)
	}
	if err=_ws.WriteMessage(websocket.TextMessage, data);err!= nil {
		// User disconnected.
		return err
	}

	return nil
}

func init(){


}