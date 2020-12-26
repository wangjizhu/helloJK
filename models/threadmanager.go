package models

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	Num_of_Threads = 6
	Num_of_Samples = Num_of_Threads * 8

	Thread_Status_Runing = iota
	Thread_Status_Paused
	Thread_Status_Stoped
	Thread_Status_Ready
)

//线程描述符
type threadDescriptor struct {
	threadName      string
	threadResources []ResourceDescriptor
	threadStatus    int
	threadParams    interface{}
}

func NewThreadDescriptor(threadname string) *threadDescriptor {
	return &threadDescriptor{
		threadName:      threadname,
		threadResources: nil,
		threadStatus:    0,
	}
}

func (th *threadDescriptor) addResource(resourceset ResourceSet) error {
	th.threadResources = append(th.threadResources, resourceset...)
	return nil
}
func (th *threadDescriptor) removeResource(resourceset ResourceSet) error {
	//查找到对应资源并充线程中删除
	for _, re := range resourceset {

		for index, value := range th.threadResources {
			if value.ResourceName == re.ResourceName && value.ResourceId == re.ResourceId {
				th.threadResources = append(th.threadResources[:index], th.threadResources[index+1:]...)
			}
		}

	}

	return nil
}

func (th *threadDescriptor) getIdbyNameFromThread(resourcename string) (int, error) {
	for _, value := range th.threadResources {
		if value.ResourceName == resourcename {
			return value.ResourceId, nil
		}
	}
	return -1, errors.New("no such resourcename in this thread")
}

var _threadTable map[string]*threadDescriptor
var wg sync.WaitGroup

//资源这一步是否需要申请和释放
type resourceApplyAndReturn struct {
	resourceName string
	apply        bool
	returnN      bool
}

func init() {
	_threadTable = make(map[string]*threadDescriptor)

	fmt.Println("threadmanager 初始化完成")

}

//要全部线程都完成才返回
func StartMultipleThreads(interval int) error {
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := StartSingleThread("111", "Sheet1")
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Duration(interval) * time.Second)

	go func() {
		defer wg.Done()
		err := StartSingleThread("222", "Sheet2")
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	return nil
}

//type resourceRecord []resourceApplyAndReturn

//在任何时候,开始一条线程,并开始执行
//调用一次开启一个线程往下走

func StartSingleThread(threadId string, sheetName string) error {
	SheetInfo := SheetsInfo[sheetName]
	//此线程已经在运行了 报错

	if _, ok := _threadTable[threadId]; ok {
		return errors.New("线程已经在运行中!")
	}
	//线程启动 注册到threadmap
	_threadTable[threadId] = NewThreadDescriptor(threadId)

	var applyresult ResourceSet

	resourcemanager, err := GetResourceManager()
	if err != nil {
		return err
	}

	fmt.Println(threadId, SheetInfo.Length)
	for i := 4; i < SheetInfo.Length; i++ {
		color.Set(color.FgBlue)
		fmt.Println(threadId, SheetInfo.Rows[i][0])
		color.Unset()

		//读入行资源
		var resourceRecord []resourceApplyAndReturn
		for j := SheetInfo.ResourceStart; j <= SheetInfo.ResourceEnd; j++ {

			rs := resourceApplyAndReturn{
				resourceName: SheetInfo.Rows[2][j],
				apply:        false,
				returnN:      false,
			}
			//	填入apply和retrun

			if SheetInfo.Rows[i][j] != "" {

				value, err := strconv.Atoi(SheetInfo.Rows[i][j])
				if err != nil {
					return err
				}

				if value/10%10 == 2 {
					rs.apply = true
				}
				if value%10 == 2 {
					rs.returnN = true
				}

			}
			resourceRecord = append(resourceRecord, rs)
		}

		//读入行参数
		params := make(map[string]string)
		for j := SheetInfo.ParamStart; j <= SheetInfo.ParamEnd; j++ {
			params[SheetInfo.Rows[2][j]] = SheetInfo.Rows[i][j]

		}

		//1.判断是否需要,并申请相应资源 申请apply为true的资源
		applyresources := ResourceSet{}
		for k := 0; k < len(resourceRecord); k++ {
			//需要申请
			if resourceRecord[k].apply == true {
				applyresources = append(applyresources,
					ResourceDescriptor{
						Status:              0,
						Enable:              false,
						ResourceName:        resourceRecord[k].resourceName,
						ResourceId:          0,
						ResourceDescription: "",
					})

			}
		}

		//有需要申请的资源
		if len(applyresources) > 0 {
			color.Set(color.FgYellow)
			fmt.Println(threadId, "applyresources:", applyresources)
			color.Unset()

			applyresult, err = resourcemanager.ApplyResource(threadId, applyresources)
			if err != nil {
				return err
			}

			//记录下来申请到的资源
			err = _threadTable[threadId].addResource(applyresult)
			if err != nil {
				return err
			}

		}
		color.Set(color.FgGreen)
		fmt.Println(threadId, "已经获得申请id", applyresult)
		color.Unset()
		//2.允许通过后 开始模拟动作 等待时间后执行完成

		duration, err := strconv.ParseFloat(SheetInfo.Rows[i][SheetInfo.Duration], 32)
		if err != nil {
			return err
		}

		err = executeAction(threadId, SheetInfo.Rows[i][0], SheetInfo.Rows[i][1], SheetInfo.Rows[i][2], params, duration)
		if err != nil {
			return err
		}

		//3.执行完成后 根据策略 归还相应资源
		returnresources := ResourceSet{}
		for k := 0; k < len(resourceRecord); k++ {
			//需要返回
			if resourceRecord[k].returnN == true {

				returnresourcedescriptor := ResourceDescriptor{
					Status:              0,
					Enable:              false,
					ResourceName:        resourceRecord[k].resourceName,
					ResourceId:          0,
					ResourceDescription: "",
				}

				//应该是从线程自己的资源里面选出来id
				id, err := _threadTable[threadId].getIdbyNameFromThread(resourceRecord[k].resourceName)
				if err != nil {
					return err
				}

				returnresourcedescriptor.ResourceId = id
				returnresources = append(returnresources, returnresourcedescriptor)
			}
		}

		if len(returnresources) > 0 {
			color.Set(color.FgYellow)
			fmt.Println(threadId, "returnresources", returnresources)
			color.Unset()

			err = resourcemanager.ReturnResource(threadId, returnresources)
			if err != nil {
				return err
			}

			//归还后 需要在线程记录中剔除次资源占用标记
			_ = _threadTable[threadId].removeResource(returnresources)

		}

	}

	//从threadmap中去除此线程
	fmt.Println(threadId, "目前线程资源管理器情况:", _threadTable[threadId].threadResources)
	delete(_threadTable, threadId)

	//从offeredmap中除去此线程

	resourcemanager.clearOfferedMapByUserName(threadId)
	fmt.Println(threadId, "目前资源管理器情况:", resourcemanager.applyingList, resourcemanager.offeredMap)

	color.Set(color.FgHiRed)
	fmt.Println(threadId, "线程执行完毕!", threadId)
	color.Unset()

	return nil
}

func getIdbyNameFromApplyResult(resourcename string, applyresult ResourceSet) (int, error) {
	for _, value := range applyresult {
		if value.ResourceName == resourcename {
			return value.ResourceId, nil
		}

	}

	fmt.Println(resourcename, applyresult)
	return -1, errors.New("no such resourcename")
}

//模拟动作执行
func executeAction(threadId string, currentStepNum string, actionName string, actionNameZH string, params interface{}, duration float64) error {

	defer timeCost(time.Now())

	//如果有必要 先做参数替换和保存 这两个动作是GP模拟器专有的 其他不用管它
	if actionName == "DP1MovePipe()" {
		//如果是这个参数 从参数库中读取 并存入线程参数
		if params.(map[string]string)["from"] == "样本装载位指定位置" {
			rm, err := GetResourceManager()
			if err != nil {
				return err
			}
			fromparam, err := rm.getIdbyNameFromBaseParams("样本装载位指定位置")
			if err != nil {
				return err
			}
			params.(map[string]string)["from"] = params.(map[string]string)["from"] + strconv.Itoa(fromparam)

			saveParamsbyIdToThread(threadId, params)

		}

	}

	if actionName == "DP1MovePipe()" {
		fmt.Println(readParamsbyIdFromThread(threadId))
		//逻辑有点怪 如果是这个参数 就需要从线程记录中读取
		if params.(map[string]string)["to"] == "样本装载位指定位置" {
			params.(map[string]string)["to"] = readParamsbyIdFromThread(threadId).(map[string]string)["from"]
		}

	}

	color.Set(color.FgHiCyan)
	fmt.Println(threadId, "动作开始:", currentStepNum, actionName, actionNameZH, params)
	color.Unset()

	_ = SendMessageThread(MessageThread{
		ThreadName: threadId,
		Resources:  MakeMessageResourceFromResourceSet(_threadTable[threadId].threadResources),
		CurrentStep: Step{
			StepDescription: "Start",
			StepOrderNum:    currentStepNum,
			StepName:        actionName,
			StepParams:      params,
		},
	})

	rand.Seed(time.Now().UnixNano())
	//执行时间自己的时间加上可能的0%-20%往上浮动
	time.Sleep(time.Duration(duration*(1+0.2*rand.Float64())*1000/20) * time.Millisecond)

	color.Set(color.FgHiMagenta)
	fmt.Println(threadId, "动作执行完毕!", actionName)
	color.Unset()

	_ = SendMessageThread(MessageThread{
		ThreadName: threadId,
		Resources:  MakeMessageResourceFromResourceSet(_threadTable[threadId].threadResources),
		CurrentStep: Step{
			StepDescription: "Finish",
			StepOrderNum:    currentStepNum,
			StepName:        actionName,
			StepParams:      params,
		},
	})

	return nil
}

//耗时统计
func timeCost(start time.Time) {
	tc := time.Since(start)
	fmt.Printf("动作耗时 = %v\n", tc)
}

func saveParamsbyIdToThread(threadid string, params interface{}) {
	_threadTable[threadid].threadParams = params
	fmt.Println("存入", params)

}

func readParamsbyIdFromThread(threadid string) interface{} {
	fmt.Println("读出", _threadTable[threadid].threadParams)
	return _threadTable[threadid].threadParams
}

func MakeMessageResourceFromResourceSet(set ResourceSet) []ResourceMessageType {
	var r []ResourceMessageType
	for _, value := range set {
		rm := ResourceMessageType{
			ResourceName:   value.ResourceName,
			ResourceAmount: strconv.Itoa(value.ResourceId),
		}
		r = append(r, rm)
	}
	return r
}
func MakeMessageResourceFromResourceMap(sm map[string]ResourceSet, samples []int) []ResourceMessageType {
	var resourcemessage []ResourceMessageType

	var keys []string
	for k := range sm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		name := k
		set := sm[k]
		amount := 0
		for _, v := range set {
			//如果为true则累加
			if v.Enable == true {
				amount++
			}

		}
		rm := ResourceMessageType{
			ResourceName:   name,
			ResourceAmount: strconv.Itoa(amount),
		}
		resourcemessage = append(resourcemessage, rm)

	}

	//map访问 每次都乱序！！！！！！！！！

	//参数资源的map也添加上去一并返回
	resourcemessage = append(resourcemessage, ResourceMessageType{
		ResourceName:   "样本装载位指定位置",
		ResourceAmount: samples,
	})

	return resourcemessage
}
