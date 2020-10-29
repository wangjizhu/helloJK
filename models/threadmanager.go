package models

import (
	"fmt"
	"strconv"
	"time"
)

//资源这一步是否需要申请和释放
type resourceApplyAndReturn struct {
	resourceName string
	apply bool
	returnN bool
}

//type resourceRecord []resourceApplyAndReturn

//在任何时候,开始一条线程,并开始执行
//调用一次开启一个线程往下走

func StartThread(threadId int)error{



	fmt.Println(len(rows))
	for i:=4;i<len(rows);i++{
		fmt.Println(i)
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





		//2.允许通过后 开始模拟动作 等待时间后执行完成

		duration,err:=strconv.Atoi(rows[i][13])
		if err != nil {
			return err
		}
		executeAction(rows[i][1],rows[i][2],params,duration,resourceRecord)



		//3.执行完成后 根据策略 归还相应资源





	}


	return nil
}



//模拟动作执行
func executeAction(actionName string,actionNameZH string,params interface{},duration int,resource []resourceApplyAndReturn){

	defer timeCost(time.Now())
	fmt.Println("动作开始:",actionName,actionNameZH,params)
	fmt.Println(resource)
	time.Sleep(time.Duration(duration/10) * time.Second)
	fmt.Println("动作执行完毕!")
}

//耗时统计
func timeCost(start time.Time){
	tc:=time.Since(start)
	fmt.Printf("动作耗时 = %v\n", tc)
}