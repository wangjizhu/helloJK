package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type ResourceMessageType struct {
	ResourceName   string
	ResourceAmount interface{}
}



type Step struct {
	StepDescription string
	StepOrderNum string
	StepName string
	StepParams interface{}
}


type MessageThread struct {
	ThreadName string
	Resources []ResourceMessageType
	CurrentStep Step
}

type MessageResource []ResourceMessageType


var _wsThread *websocket.Conn
var _wsResource *websocket.Conn
var _wsLock sync.Mutex

//var _wsThreadSendBuffer chan MessageThread
//var _wsResourceSendBuffer chan MessageResource

//链接一次 记下来
func SetWsThread(ws *websocket.Conn){
	_wsThread =ws
}
func GetWsThread()*websocket.Conn{
	return _wsThread
}

func CloseWsThread(){
	if _wsThread !=nil{
		_wsThread.Close()
		fmt.Println("已经关闭当前ws")
	}
}

func SendMessageThread(m MessageThread)error{
	data,err:=json.Marshal(&m)
	if err!=nil{
		panic(err)
	}

	if _wsThread==nil{
		return errors.New("no _wsThread")
	}

	//并发发送会带来问题 故需要枷锁来解决

	_wsLock.Lock()

	if err= _wsThread.WriteMessage(websocket.TextMessage, data);err!= nil {
		// User disconnected.
		return err
	}
	_wsLock.Unlock()

	return nil
}

//---------------------------------------------------------------------
//链接一次 记下来
func SetWsResource(ws *websocket.Conn){
	_wsResource =ws
}
func GetWsResource()*websocket.Conn{
	return _wsResource
}

func CloseWsResource(){
	if _wsResource !=nil{
		_wsResource.Close()
		fmt.Println("已经关闭当前ws")
	}
}





func init(){
	//_wsThreadSendBuffer:=make(chan MessageThread)
	//_wsResourceSendBuffer:=make(chan MessageResource)



}