package models

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
	"sync"
	"time"
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

	rows, err = _f.GetRows("Sheet1")
	if err != nil {
		panic(err)
	}

	fmt.Println(len(rows))

	_r=NewResourceManager()
	err=_r.Init()
	if err != nil {
		panic(err)
	}

	fmt.Println(_r.resources)

	err = StartThread(1)
	if err != nil {
		panic(err)
	}



}

//资源描述 申请和响应都用同样结构 根据实际情况使用不同字段
type ResourceDescriptor struct {
	//查询是返回状态
	Status int
	Enable bool

	//资源名 资源id 资源备注
	ResourceName string
	ResourceId int
	ResourceDescription string
}

type ResourceSet [] ResourceDescriptor


type applyItem struct{
	userName          string
	applyResourceSet ResourceSet    //挂一系列需要申请的资源
	resultChan       chan ResourceSet

}



type ResourceManager struct {
	//这个锁用来保证map多线程查询和更改时候的安全性 官方解决方案
	sync.RWMutex
	listLock sync.Mutex
	queryLock sync.Mutex

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

func GetResourceManager()(*ResourceManager,error){
	if _r!=nil{
		return _r,nil
	}else{
		return nil,errors.New("no rsmanger at all")
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
				ResourceId:          j,
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



//向资源管理器申请某些资源 当全部都可以被使用时候返回 否则阻塞 可能会超时
func (r *ResourceManager)ApplyResource(username string ,resources ResourceSet) (ResourceSet,error){
	fmt.Println("提出申请",username,resources)

	var getresourceset ResourceSet
	//测试 一旦申请 直接满足回复 用于单线程测试
	applyitem:=applyItem{
		userName:         username,
		applyResourceSet: resources,
		resultChan:       make(chan ResourceSet),
	}

	//注册到申请列表
	r.listLock.Lock()
	r.applyingList.PushBack(applyitem)
	r.listLock.Unlock()
	fmt.Println("推进去一个申请")

	//阻塞等待 可能会超时
	select{
		case getresourceset = <-applyitem.resultChan:
			fmt.Println("成功获得资源",username,resources)
			return getresourceset,nil

		case <-time.After(100 * time.Second):
			fmt.Println("申请等待超时!",username,resources)
			return nil,errors.New("timeout")
	}

}


//归还资源
func (r *ResourceManager)ReturnResource(username string ,resources ResourceSet)error{
	fmt.Println("开始归还",username,resources)

	return nil
}

//核心!通过查找applyList 查找是否能够满足用户需求 并返回用户申请 ,,,每秒调用一次
//若是多资源申请 要全部资源都可用情况下 才为可用 并返回全部资源到对应applyItem的resultChan里面
func (r *ResourceManager)CheckList()error {
	var next *list.Element
	for e := r.applyingList.Front(); e != nil; e = next {
		// do something with e.Value 天坑这里 先把next存下来 万一删除了就找不到下一个了
		next = e.Next()

		//检查这个节点里面的资源是否全部可用
		available,err:= r.allResourcesAvailable(e.Value.(applyItem))
		if err != nil {
			return err
		}

		//全部资源可用
		if available==true{




		//	不是全部资源可用
		}else{


		}



	}

	return nil
}

func (r *ResourceManager)allResourcesAvailable(applyitem applyItem)(bool,error){
	//查询时锁住 其实没必要 因为一次只能一个查询
	//r.queryLock.Lock()
	allright:=true

	for i:=0;i<len(applyitem.applyResourceSet);i++{
		name:=applyitem.applyResourceSet[i].ResourceName
		id:=applyitem.applyResourceSet[i].ResourceId
		if r.resources[name][id].Enable==true{


		}else{

			allright=false
			break

		}


	}

	//都有
	if allright==true{

		return true,nil
	//	有一个没有
	}else{

		return false,nil

	}



	//查完了 放开
	//r.queryLock.Unlock()


}