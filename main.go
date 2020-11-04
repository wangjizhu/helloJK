package main

import (
	"github.com/astaxie/beego"
	_ "helloprecision/routers"
)

func main() {

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		beego.BConfig.WebConfig.StaticDir["/static"] = "static"
		//beego.BConfig.WebConfig.TemplateLeft = "[["
		//beego.BConfig.WebConfig.TemplateRight = "]]"
	}


	go beego.Run()

	//go func() {
	//	err:=models.StartMultipleThreads(models.Num_of_Threads,80)
	//	if err != nil {
	//		panic("fatal error")
	//	}
	//}()



	select {

	}

}

