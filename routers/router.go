package routers

import (
	"helloprecision/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})

	ns1 := beego.NewNamespace("/v1",
		beego.NSNamespace("/GP", beego.NSInclude(&controllers.UserInterfaceController{})))


	beego.AddNamespace(ns1)


}
