package raptorQ

import (
//	"fmt"
//	"fmt"
)

type Symbol struct {
	data []byte
	nbytes int    // number of data
	sbn int       // source block number
	esi int       // block number of internal
}



func (ss *Symbol) xxor(s *Symbol) {
	if ss.data==nil{
		ss.data=make([]byte,ss.nbytes)
	}
    if ss.nbytes != s.nbytes {
		return 
	}
	for i:=0;i<ss.nbytes; i++{
		ss.data[i] = ss.data[i] ^ s.data[i]
	}

}

func (ss *Symbol) init(size int) {
	if ss.nbytes != size{
		ss.nbytes = size
		ss.data = make([]byte,size)
	}
	ss.data=make([]byte,size)
}


func (ss *Symbol) fillData(src []byte,size int){
	if ss.nbytes != size {
		ss.nbytes = size
		ss.data = make([]byte,size)
	}
		ss.data=src
}

func (ss *Symbol) muladd(s Symbol,u byte){

	if ss.data==nil{
		ss.data=make([]byte,ss.nbytes)
	}
	if ss.nbytes != s.nbytes{
	    return
	}

	for i:=0;i<ss.nbytes; i++{
		if int(u)==1{
		     ss.data[i] ^= s.data[i]
		}else{
			 ss.data[i] ^= octmul(s.data[i],u)
		}
	}
}

func (ss *Symbol) div(u byte){
	if ss.data==nil{
		ss.data=make([]byte,ss.nbytes)
	}
	for i:=0;i<ss.nbytes; i++{
		ss.data[i] = octdiv(ss.data[i],u)
	}
}


func (ss *Symbol) mul(u byte){
	if ss.data==nil{
		ss.data=make([]byte,ss.nbytes)
	}
	for i:=0;i<ss.nbytes; i++{
		ss.data[i] = octmul(ss.data[i],u)
	}
}


