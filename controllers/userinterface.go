package controllers

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/adapter"
	"helloprecision/models"
	"net/http"
)

type GPController struct {
	beego.Controller
}

// @Title SetResourceSample
// @Description 设置原始样品试管内容
// @Param	body	body 	body 	true	"SetResourceSample"
// @Success 200 ok
// @Failure 403 no
// @router /SetResourceSample/ [post]
func (u *GPController) SetResourceSample(){
	fmt.Println("aaaaaaa")

	var s []int
	err:=json.Unmarshal(u.Ctx.Input.RequestBody, &s)
	fmt.Println(s)
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
// @router /StartSingleThread/ [get]
func (u *GPController) StartSingleThread(){
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
// @Param	interval	query   int	 true ">=80"
// @Success 200 ok
// @Failure 403 no
// @router /StartMultipleThreads/ [get]
func (u *GPController) StartMultipleThreads(){
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

// @Title GetLengthOfThread
// @Description 获取线程总长度
// @Success 200 ok
// @Failure 403 no
// @router /GetLengthOfThread/ [get]
func (u *GPController) GetLengthOfThread(){
	u.Data["json"]=models.BookInfo.Length
	u.ServeJSON()
	return
}


// @Title BorrowSampleShelf
// @Description BorrowSampleShelf
// @Success 200 ok
// @Failure 403 no
// @router /BorrowSampleShelf/ [get]
func (u *GPController) BorrowSampleShelf(){
	applyresources:=models.ResourceSet{}
	applyresources=append(applyresources,models.ResourceDescriptor{
		Status:              0,
		Enable:              false,
		ResourceName:        "SampleShelf",
		ResourceId:          0,
		ResourceDescription: "",
	})

	resourcemanager,err:=models.GetResourceManager()
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusInternalServerError)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return
	}

	_,err=resourcemanager.ApplyResource("",applyresources)


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

// @Title ReturnSampleShelf
// @Description ReturnSampleShelf
// @Success 200 ok
// @Failure 403 no
// @router /ReturnSampleShelf/ [get]
func (u *GPController) ReturnSampleShelf(){
	applyresources:=models.ResourceSet{}
	applyresources=append(applyresources,models.ResourceDescriptor{
		Status:              0,
		Enable:              false,
		ResourceName:        "SampleShelf",
		ResourceId:          0,
		ResourceDescription: "",
	})

	resourcemanager,err:=models.GetResourceManager()
	if err != nil {
		u.Ctx.Output.SetStatus(http.StatusInternalServerError)
		u.Data["json"]=err.Error()
		u.ServeJSON()
		return
	}

	err=resourcemanager.ReturnResource("user",applyresources)

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