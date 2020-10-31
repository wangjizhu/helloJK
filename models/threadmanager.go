package models

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"github.com/fatih/color"
)


const (
	Thread_Status_Runing=iota
	Thread_Status_Paused
	Thread_Status_Stoped
	Thread_Status_Ready

)

//线程描述符
type threadDescriptor struct{
	threadName string
	threadResources []ResourceDescriptor
	threadStatus int
}

func NewThreadDescriptor(threadname string)*threadDescriptor{
	return &threadDescriptor{
		threadName:      threadname,
		threadResources: nil,
		threadStatus:    0,
	}
}

func (th* threadDescriptor) addResource(resourceset ResourceSet) error{
	th.threadResources = append(th.threadResources, resourceset...)
	return nil
}
func (th* threadDescriptor) removeResource(resourceset ResourceSet) error{
	//查找到对应资源并充线程中删除
	for _,re:=range resourceset{

		for index,value:=range th.threadResources{
			if value.ResourceName==re.ResourceName && value.ResourceId==re.ResourceId{
				th.threadResources = append(th.threadResources[:index],th.threadResources[index+1:]... )
			}
		}

	}

	return nil
}

func (th* threadDescriptor) getIdbyNameFromThread(resourcename string) (int,error){
	for _,value:=range th.threadResources{
		if value.ResourceName==resourcename {
			return value.ResourceId,nil
		}
	}
	return -1,errors.New("no such resourcename in this thread")
}

var  _threadTable map[string] *threadDescriptor


//资源这一步是否需要申请和释放
type resourceApplyAndReturn struct {
	resourceName string
	apply bool
	returnN bool
}

func init(){
	_threadTable=make(map[string] *threadDescriptor)

	fmt.Println("threadmanager 初始化完成")

	for i:=1;i<=100;i++{
		go func() {
			err := StartThread(fmt.Sprintf("%d%d%d",i,i,i))
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(70*time.Second)
	}

}






//type resourceRecord []resourceApplyAndReturn

//在任何时候,开始一条线程,并开始执行
//调用一次开启一个线程往下走

func StartThread(threadId string)error{

	//线程启动 注册到threadmap
	_threadTable[threadId]=NewThreadDescriptor(threadId)



	var applyresult ResourceSet

	resourcemanager,err:=GetResourceManager()
	if err != nil {
		return err
	}


	fmt.Println(threadId,len(rows))
	for i:=4;i<len(rows);i++{
		color.Set(color.FgBlue)
		fmt.Println(threadId,rows[i][0])
		color.Unset()

		//读入行资源
		var resourceRecord []resourceApplyAndReturn
		for j:=4;j<=11;j++{


			rs:=resourceApplyAndReturn{
				resourceName: rows[2][j],
				apply:        false,
				returnN:      false,
			}
		//	填入apply和retrun

			if rows[i][j]!=""{

				value, err := strconv.Atoi(rows[i][j])
				if err != nil {
					return err
				}

				if value/10%10==2{
					rs.apply=true
				}
				if value%10==2{
					rs.returnN=true
				}

			}
			resourceRecord=append(resourceRecord, rs)
		}

		//读入行参数
		params :=make(map[string]string)
		for j:=14;j<=23;j++{
			params[rows[2][j]]=rows[i][j]

		}

		//1.判断是否需要,并申请相应资源 申请apply为true的资源
		applyresources:=ResourceSet{}
		for k:=0;k<len(resourceRecord);k++{
			//需要申请
			if resourceRecord[k].apply==true{
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
		if len(applyresources)>0{
			color.Set(color.FgYellow)
			fmt.Println(threadId,"applyresources:",applyresources)
			color.Unset()


			applyresult,err=resourcemanager.ApplyResource(threadId,applyresources)
			if err != nil {
				return err
			}

			//记录下来申请到的资源
			err=_threadTable[threadId].addResource(applyresult)
			if err != nil {
				return err
			}

		}
		color.Set(color.FgGreen)
		fmt.Println(threadId,"已经获得申请id",applyresult)
		color.Unset()
		//2.允许通过后 开始模拟动作 等待时间后执行完成

		duration,err:=strconv.ParseFloat(rows[i][13], 32)
		if err != nil {
			return err
		}
		executeAction(threadId,rows[i][1],rows[i][2],params,duration)



		//3.执行完成后 根据策略 归还相应资源
		returnresources:=ResourceSet{}
		for k:=0;k<len(resourceRecord);k++{
			//需要返回
			if resourceRecord[k].returnN==true{

				returnresourcedescriptor:=ResourceDescriptor{
					Status:              0,
					Enable:              false,
					ResourceName:        resourceRecord[k].resourceName,
					ResourceId:          0,
					ResourceDescription: "",
				}

				//应该是从线程自己的资源里面选出来id
				id,err:=_threadTable[threadId].getIdbyNameFromThread(resourceRecord[k].resourceName)
				if err != nil {
					return err
				}

				returnresourcedescriptor.ResourceId=id
				returnresources = append(returnresources, returnresourcedescriptor)
			}
		}

		if len(returnresources)>0{
			color.Set(color.FgYellow)
			fmt.Println(threadId,"returnresources",returnresources)
			color.Unset()

			err=resourcemanager.ReturnResource(threadId,returnresources)
			if err != nil {
				return err
			}

			//归还后 需要在线程记录中剔除次资源占用标记
			_threadTable[threadId].removeResource(returnresources)

		}


	}

	//从threadmap中去除此线程
	fmt.Println(threadId,"目前线程资源管理器情况:",_threadTable[threadId].threadResources)
	delete(_threadTable,threadId)

	//从offeredmap中除去此线程
	fmt.Println(threadId,"目前资源管理器情况:",resourcemanager.applyingList,resourcemanager.offeredMap)
	resourcemanager.clearOfferedMapByUserName(threadId)

	color.Set(color.FgCyan)
	fmt.Println(threadId,"线程执行完毕!",threadId)
	color.Unset()



	return nil
}

func getIdbyNameFromApplyResult(resourcename string,applyresult ResourceSet) (int,error){
	for _,value := range applyresult{
		if value.ResourceName==resourcename{
			return value.ResourceId,nil
		}

	}

	fmt.Println(resourcename,applyresult)
	return -1,errors.New("no such resourcename")
}


//模拟动作执行
func executeAction(threadId string,actionName string,actionNameZH string,params interface{},duration float64){

	defer timeCost(time.Now())

	color.Set(color.FgHiCyan)
	fmt.Println(threadId,"动作开始:",actionName,actionNameZH,params)
	color.Unset()

	rand.Seed(time.Now().UnixNano())
	//执行时间自己的时间加上可能的0%-20%往上浮动
	time.Sleep(time.Duration(duration*(1+0.2*rand.Float64())*1000/20) * time.Millisecond)


	color.Set(color.FgHiMagenta)
	fmt.Println(threadId,"动作执行完毕!",actionName)
	color.Unset()
}

//耗时统计
func timeCost(start time.Time){
	tc:=time.Since(start)
	fmt.Printf("动作耗时 = %v\n", tc)
}