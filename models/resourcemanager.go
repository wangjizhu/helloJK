package models

import (
	"container/list"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
	"sync"
)

var _r * ResourceManager
var _f *excelize.File
var rows [][]string


func init(){
	var err error
	_f,err=excelize.OpenFile("./Book1.xlsx")
	if err != nil {
		panic(err)
	}
	fmt.Println("aaa")

	rows, err = _f.GetRows("Sheet2")
	if err != nil {
		panic(err)
	}

	_r=NewResourceManager()
	err=_r.Init()
	if err != nil {
		panic(err)
	}

	fmt.Println(_r.resources)



}

//资源描述 申请和响应都用同样结构 根据实际情况使用不同字段
type ResourceDescriptor struct {
	//查询是返回状态
	Status int
	Enable bool

	//资源名 资源id 资源备注
	ResourceName string
	ResourceId string
	ResourceDescription string
}

type ResourceSet [] ResourceDescriptor


type ResourceManager struct {
	//这个锁用来保证map多线程查询和更改时候的安全性 官方解决方案
	sync.RWMutex

	//全部的资源纪录
	resources map[string] ResourceSet

	//正在申请中的list 需要满足 谁想要接什么
	applyingList *list.List

	//已经被借走的资源列表 谁借了什么
	offeredList *list.List

}


func NewResourceManager()*ResourceManager{
	return &ResourceManager{
		RWMutex:      sync.RWMutex{},
		resources:    nil,
		applyingList: nil,
		offeredList:  nil,
	}

}

//根据配置 初始化资源管理器 主要是置其中各种资源
func (r *ResourceManager)Init()error{

	r.resources=make(map[string] ResourceSet)


	for i:=4;i<=11;i++{
		//挂资源名
		rname:=rows[2][i]
		r.resources[rname]=ResourceSet{}
		fmt.Println(rname)

		//挂资源内容
		amount, err := strconv.Atoi(rows[3][i])
		if err != nil {
			return err
		}
		fmt.Println(amount)

		for j:=0;j<amount;j++{
			resource:=ResourceDescriptor{
				Status:              0,
				Enable:              true,
				ResourceName:        rname,
				ResourceId:          strconv.Itoa(j),
				ResourceDescription: "",
			}
			r.resources[rname]=append(r.resources[rname],resource)
		}

	}

	return nil
}

//向资源管理器注册一个新的资源
func (r *ResourceManager) registerResource(resourcename string,recourceamount int)error {


	return nil
}



//向资源管理器申请某些资源 当全部都可以被使用时候返回 否则阻塞
func (r *ResourceManager)Apply(username string ,resources ResourceSet) error{


	return nil
}

func (r *ResourceManager)Return(username string ,resources ResourceSet)error{

	return nil
}
