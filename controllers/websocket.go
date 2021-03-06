// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	"helloprecision/models"
	//"encoding/json"
	"net/http"

	beego "github.com/beego/beego/v2/adapter"
	"github.com/gorilla/websocket"
	//"github.com/beego/samples/WebIM/models"
	"fmt"
	//"time"

	//"encoding/json"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	beego.Controller
}


// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Join() {
	fmt.Println("Thread ws链接成功")

	//关闭已有的
	defer models.CloseWsThread()

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	//记录下新的
	models.SetWsThread(ws)

	//检测链接错误使用 若无法读取 则断开ws
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			return
		}

	}





}




func (this *WebSocketController) ResourceStatus() {
	fmt.Println("Resource ws链接成功")

	//关闭已有的
	defer models.CloseWsResource()

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	//记录下新的
	models.SetWsResource(ws)

	//检测链接错误使用 若无法读取 则断开ws
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			return
		}

	}




}