package models

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"strconv"
	"sync"
	"time"
)





var _r * ResourceManager
var _f *excelize.File
var rows [][]string
var LengthOfRows int


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

	LengthOfRows=len(rows)
	fmt.Println(LengthOfRows)

	_r=NewResourceManager()
	err=_r.Init()
	if err != nil {
		panic(err)
	}

	fmt.Println(_r.resources)


	//一秒一次轮询用于list
	go func() {
		for{
			_ = _r.PollingList()
			time.Sleep(1000 * time.Millisecond)
		}
	}()



}

//资源描述符 申请和响应都用同样结构 根据实际情况使用不同字段
type ResourceDescriptor struct {
	//查询时返回状态
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
	baseParams map[string] []int

	//正在申请中的list 需要满足 谁想要接什么
	applyingList *list.List

	//已经被借走的资源列表 谁借了什么
	offeredMap map[string] [] ResourceDescriptor

}


func NewResourceManager()*ResourceManager{
	return &ResourceManager{
		RWMutex:      sync.RWMutex{},
		resources:    nil,
		applyingList: nil,
		offeredMap:  nil,
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

	//初始化两个list
	r.applyingList=list.New()


	r.resources=make(map[string] ResourceSet)
	r.offeredMap=make(map[string] [] ResourceDescriptor)

	r.baseParams=make(map[string] []int)
	for s:=0;s<Num_of_Samples;s++{
		r.baseParams["样本装载位指定位置"] = append(r.baseParams["样本装载位指定位置"], s+11)
	}


	for i:=4;i<=13;i++{
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

	//每秒返回资源管理器中资源数量
	go func() {

		for{
			//fmt.Println(r.resources)
			err:=SendMessageResource(MakeMessageResourceFromResourceMap(r.resources,r.baseParams["样本装载位指定位置"]))
			if err != nil {

			}
			time.Sleep(1*time.Second)
		}

	}()



	return nil
}

//向资源管理器注册一个新的资源
func (r *ResourceManager) registerResource(resourcename string,recourceamount int)error {

	return nil
}



//向资源管理器申请某些资源 可以一次申请多个 当全部都可以被使用时候返回 否则阻塞 可能会超时
func (r *ResourceManager)ApplyResource(username string ,resources ResourceSet) (ResourceSet,error){
	//不需要申请 就不要瞎参和
	if len(resources)==0{
		return nil,errors.New("empty resource apply")
	}



	var getresourceset ResourceSet
	//测试 一旦申请 直接满足回复 用于单线程测试
	applyitem:=applyItem{
		userName:         username,
		applyResourceSet: resources,
		resultChan:       make(chan ResourceSet),
	}

	fmt.Println("提出申请",applyitem)

	//注册到申请列表
	r.listLock.Lock()
	r.applyingList.PushBack(applyitem)
	r.listLock.Unlock()
	fmt.Println("推进去一个申请")

	//阻塞等待 可能会超时
	select{
		case getresourceset = <-applyitem.resultChan:
			color.Set(color.BgGreen,color.FgBlack)
			fmt.Println("成功获得资源",username,resources)
			color.Unset()
			return getresourceset,nil

		case <-time.After(60 * time.Second):
			color.Set(color.BgRed,color.FgBlack)
			fmt.Println("申请等待超时!")
			fmt.Println(username,resources)
			color.Unset()
			return nil,errors.New("timeout")
	}

}


//归还资源 强调 归还的资源需要携带id!!!
func (r *ResourceManager)ReturnResource(username string ,resources ResourceSet)error{
	fmt.Println("开始归还",username,resources)

	for i:=0;i<len(resources);i++{
		resourcename:=resources[i].ResourceName
		id:=resources[i].ResourceId

		//资源管理器中还回去
		r.resources[resourcename][id].Enable=true
		//offerdMap中去掉
		err:=r.removeResourceFromOfferedMap(username,resourcename,id)
		if err != nil {
			return err
		}
	}

	fmt.Println("已归还",username,resources)
	return nil
}

func (r *ResourceManager) removeResourceFromOfferedMap(username string,resourcename string,id int)error{
	for i:=0;i<len(r.offeredMap[username]);i++{
		//找到后删除
		if r.offeredMap[username][i].ResourceName==resourcename && r.offeredMap[username][i].ResourceId==id{
			r.offeredMap[username]=append(r.offeredMap[username][:i], r.offeredMap[username][i+1:]...)
		}
	}
	return nil
}

//核心!通过查找applyList 查找是否能够满足用户需求 并返回用户申请 ,,,每秒调用一次
//若是多资源申请 要全部资源都可用情况下 才为可用 并返回全部资源到对应applyItem的resultChan里面
func (r *ResourceManager)PollingList()error {

	//fmt.Println("进入polling...")
	var next *list.Element
	for e := r.applyingList.Front(); e != nil; e = next {
		// do something with e.Value 天坑这里 先把next存下来 万一删除了就找不到下一个了
		next = e.Next()
		applyitem:=e.Value.(applyItem)
		fmt.Println("轮询中资源需求",applyitem)
		//检查这个节点里面的资源是否全部可用
		available,checkresult,err:= r.checkAllResourcesAvailableFromOneApplyItem(applyitem)
		if err != nil {
			return err
		}

		//全部资源可用
		if available==true{
			//去资源管理器记录下来
			err:=r.getAllResourcesAvailableFromResult(checkresult)
			if err != nil {
				return err
			}

			//申请已经得到满足 从申请列表中去除
			r.applyingList.Remove(e)

			//表示申请成功 挂入offeredMap 记录下被谁取走
			for i:=0;i<len(checkresult.applyResourceSet);i++{
				//如果不存在则会新建
				r.offeredMap[applyitem.userName]=append(r.offeredMap[applyitem.userName], checkresult.applyResourceSet[i])

			}


			fmt.Println(checkresult.userName,"获得资源",checkresult.applyResourceSet)

			//结果用通道返回
			e.Value.(applyItem).resultChan <- checkresult.applyResourceSet


		//	不是全部资源可用 那就只有算求
		}else{


		}



	}

	return nil
}


//检查节点内部的是否全部资源可用 返回是否可用,可用的资源,和错误
func (r *ResourceManager)checkAllResourcesAvailableFromOneApplyItem(applyitem applyItem)(bool,applyItem,error){
	//查询时锁住 其实没必要 因为一次只能一个查询
	//r.queryLock.Lock()
	allright:=true

	for i:=0;i<len(applyitem.applyResourceSet);i++{
		name:=applyitem.applyResourceSet[i].ResourceName

		//id是不知道的
		affect,id,err:=r.checkResourceEnableFromResourceMap(name)
		if err != nil {
			return false,applyitem,err
		}

		//找到了资源
		if affect==true{
			//id记录下来
			applyitem.applyResourceSet[i].ResourceId=id

		}else{

			allright=false
			break

		}

	}

	//都有
	if allright==true{

		return true,applyitem,nil
	//至少有一个没有
	}else{

		return false,applyitem,nil

	}



	//查完了 放开
	//r.queryLock.Unlock()

}

//检查并返回具体资源可用的那个Id 没找到则返回false
func (r *ResourceManager) checkResourceEnableFromResourceMap(resourcename string)(bool,int,error){

	for i:=0;i<len(r.resources[resourcename]);i++{

		if r.resources[resourcename][i].Enable==true{
			return	true, i, nil
		}

	}

	//没有空闲资源
	return false, 0, nil
}





//通过对资源管理器资源标记,来把资源分配给用户 在先调用可用判断后才能分配
func (r *ResourceManager)getAllResourcesAvailableFromResult(checkresult applyItem)error{
	for i:=0;i<len(checkresult.applyResourceSet);i++{
		name:=checkresult.applyResourceSet[i].ResourceName
		id:=checkresult.applyResourceSet[i].ResourceId
		//登记处取下来
		r.resources[name][id].Enable=false
	}
	return nil
}

//清理offeredmap 把对应线程的内容删除掉
func (r *ResourceManager)clearOfferedMapByUserName(username string){
	delete(r.offeredMap,username)
}


func (r *ResourceManager)getIdbyNameFromBaseParams(paramname string)(int,error){

	if len(r.baseParams[paramname])>0{

		value:=r.baseParams[paramname][0]
		r.baseParams[paramname] = append(r.baseParams[paramname][:0], r.baseParams[paramname][1:]...)

		return value,nil

	}

	return -1,errors.New(paramname+"is empty!")
}

func SetResourceSample(s []int)error{
	r,err:=GetResourceManager()
	if err != nil {
		return err
	}
	r.baseParams["样本装载位指定位置"]=s
	return nil
}


func SendMessageResource(m MessageResource)error{
	data,err:=json.Marshal(&m)
	if err!=nil{
		panic(err)
	}
	ws:=GetWsResource()

	if ws==nil{
		return errors.New("no _wsResource")
	}

	if err= ws.WriteMessage(websocket.TextMessage, data);err!= nil {
		// User disconnected.
		return err
	}

	return nil
}
