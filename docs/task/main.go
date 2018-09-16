package main
import (
	"fmt"
	"os"
	"io/ioutil"
	"bytes"
	"strconv"
//	"time"
	
)

type matrix int
const (
	matrix1=4
	matrix2=5
	police=1
	car=2
)

func main(){
	//start:=time.Now()
	f,_:=os.Open("input2.txt")
	buf,_:=ioutil.ReadAll(f)
	ds:=&DataSet{}
	ds.Init(buf)
	ma:=setup(ds.matrix)
	for i:=0;i<ds.car*12;i++{
		ma.Calculate(ds.data[i].x,ds.data[i].y)
	}

	number:=0
    if ds.police==1{
		for i:=0;i<ds.matrix;i++{
			for ii:=0;ii<ds.matrix;ii++{
				if ma.M[i][ii]>number{
					number=ma.M[i][ii]
				}
			}
		}
	} else {

	}
//	end:=time.Now()

  pp:= Plicepoint(5,1)

 for i:=0;i<5;i++{
	 fmt.Println(pp[i])
 }
    fmt.Println(number)


}

type Matrix struct{
	M [][]int
}

type Point struct{
	flag byte
	x int
	y int
}

type DataSet struct{ 
	matrix int
	police  int
	car     int
	data    []Point
}


func (d *DataSet) Init(set []byte){
   len:=bytes.Count(set,[]byte{13,10})
   for i:=0;i<len;i++{
	index:=bytes.IndexAny(set,"\r\n")
	if i==0{
		result,_:=strconv.Atoi(string(set[0:index]))
		d.matrix=result
	}else if i==1{
		result,_:=strconv.Atoi(string(set[0:index]))
		d.police=result
	}else if i==2{
		result,_:=strconv.Atoi(string(set[0:index]))
		d.car=result
	}else {
		dd:=set[0:index]
		x,_:=strconv.Atoi(string(dd[0:1]))
		y,_:=strconv.Atoi(string(dd[2:]))
		d.data=append(d.data,Point{x:x,y:y})
	}
	copy(set[0:], set[index+2:])
   }
}


func setup(n int) *Matrix{
	 matrix:=make([][]int,n)
	 for i:=0;i<n;i++{
		 matrix[i]=make([]int,n)
	 }
	 return &Matrix{matrix}
}

func (m *Matrix) Calculate(x,y int){
     m.M[x][y]+=1;
}

type possibility struct{
   moniter Point
   block   []Point
} 

func Plicepoint(n,p int) [][]Point{
	matriax:=make([][]Point,n)
    for i:=0;i<n;i++{
		matriax[i]=make([]Point,n)
	}
	moniter:=Point{x:2,y:1}
	for i:=0;i<n;i++{
		for ii:=0;ii<n;ii++{
			if i==moniter.x || ii==moniter.y ||(i==(ii-moniter.y))|| (i==(moniter.y-ii)){
			}else{
				matriax[i][ii].flag=1
				matriax[i][ii].x=i
				matriax[i][ii].y=ii
			}
		}
	}
	matriax[moniter.x][moniter.y].flag=1
	matriax[moniter.x][moniter.y].x=moniter.x
	matriax[moniter.y][moniter.y].y=moniter.y
	
	

	return matriax
}