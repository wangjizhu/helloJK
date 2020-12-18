package routers

import (
	"helloprecision/controllers"
	beego "github.com/beego/beego/v2/adapter"
)

func init() {

	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")
	beego.Router("/ws/resource", &controllers.WebSocketController{}, "get:ResourceStatus")

    beego.Router("/", &controllers.MainController{})

	ns1 := beego.NewNamespace("/v1",
		beego.NSNamespace("/GP", beego.NSInclude(&controllers.GPController{})))


	beego.AddNamespace(ns1)


}
