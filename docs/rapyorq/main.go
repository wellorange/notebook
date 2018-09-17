package main

import (
	"./raptorQ"
	 //"fmt"
)

func main(){
	var  comm raptorQ.Common
	message := [][]byte{
		{97,98,99,100},
		{97,98,99,100},
		{97,98,99,100},
		{97,98,99,100},
		{97,98,99,100},
		{97,98,99,100},
		{97,98,99,100},
		{97,98,99,100}}
	K := 8
	T := 4
	lossrate := 0.1
	overhead := (int)((float64(K)*lossrate + 10) / (1 - lossrate))
	comm.Prepare(K,K,T)

	// 编码实现
	comm.Code(message,overhead,nil)



	// 生成中间符号
	inter:=comm.Generate_intermediates()
	for i:=0;i<len(inter);i++{
		//fmt.Println(inter[i])
	}
	//fmt.Println()
	//fmt.Println()

	// 生成修正符号
	repair:=comm.Generate_repairs(overhead)

	for i:=0;i<len(repair);i++{
		//fmt.Println(repair[i])
	}

}
