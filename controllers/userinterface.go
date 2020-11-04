package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"helloprecision/models"
	"net/http"
)

type UserInterfaceController struct {
	beego.Controller
}

// @Title SetResourceSample
// @Description 设置原始样品试管内容
// @Param	body	body 	body 	true	"SetResourceSample"
// @Success 200 ok
// @Failure 403 no
// @router /SetResourceSample/ [post]
func (u *UserInterfaceController) SetResourceSample(){
	var s []int
	err:=json.Unmarshal(u.Ctx.Input.RequestBody, &s)
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusBadRequest)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return
	}

	err=models.SetResourceSample(s)
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusInternalServerError)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return

	}

	u.Data["json"]="ok"
	u.ServeJSON()
	return
}

// @Title startSingleThread
// @Description 输入线程号 开始单线程
// @Param	threadId	query   string	 true "something"
// @Success 200 ok
// @Failure 403 no
// @router /startSingleThread/ [get]
func (u *UserInterfaceController) StartSingleThread(){
	threadId:=u.GetString("threadId")

	err:=models.StartSingleThread(threadId)
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusInternalServerError)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return
	}

	u.Data["json"]="ok"
	u.ServeJSON()
	return
}

// @Title StartMultipleThreads
// @Description 开始多条线程 输入线程数量 和线程开始时间间隔 将自动开始 直到全部线程完成后 调用结束
// @Param	numOfThread	query   int	 true "something"
// @Param	interval	query   int	 true "something"
// @Success 200 ok
// @Failure 403 no
// @router /StartMultipleThreads/ [get]
func (u *UserInterfaceController) StartMultipleThreads(){
	numOfThread,err:=u.GetInt("numOfThread")
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusBadRequest)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return
	}
	interval,err:=u.GetInt("interval")
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusBadRequest)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return
	}

	err=models.StartMultipleThreads(numOfThread,interval)
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusInternalServerError)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return
	}

	u.Data["json"]="ok"
	u.ServeJSON()
	return
}