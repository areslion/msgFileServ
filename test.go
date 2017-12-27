package main
import "fmt"
import "unsafe"


func main(){
	fmt.Println()
	//tstGetMaxMin()
	tstPointer();
}

func tstPointer(){
	a :=10
	fmt.Println(a,&a)
}

func tstGetMaxMin(){
	lst :=[]int{1,2,-1,456,100,-10}
	fmt.Println(tstArrayMaxMin(lst))
}

func tstArrayMaxMin(lst []int)(max int,min int){
	max,min = lst[0],lst[0]
	for index:=range(lst){
		if(lst[index]>max){
			max = lst[index]
		}
		if(lst[index]<min){
			min = lst[index]
		}
	}

	return max,min
}

func tstArray(lst [] int)(max int ,min int){
	min = lst[0]
	max = lst[0]
	
	for index := 0; index < len(lst); index++ {
		if(lst[index]<min) {
			min=lst[index]
		}
		if(lst[index]>max) {
			max=lst[index]
		}
	}



	return max,min
}

func tstMaxMin(a int ,b int)(max int ,min int){
		if(a>b) {
			max=a;min=b
		} else {
			max=b;min=a
		}
	
	return max,min
}

func tstConst(){
	const (
		PI = "3.1415926"
		size = len(PI)
		sizex = unsafe.Sizeof(PI)
	)
	const (
		a1 = iota
		a2
		a3
		a4 = "foure"
		a5
		a6 = 66
		a7
		a8 = iota
		a9
		a10
	)

	fmt.Println(PI,size,sizex)
	fmt.Println(a1,a2,a3,a4,a5,a6,a7,a8,a9,a10)
}

func tstVar(){
	var a,b,c = 1,1.1,"hello"
	fmt.Println(a,b,c)

	a1,b1,c1 :=2,2.2,"second"
	fmt.Println(a1,b1,c1);
}
func tstHello(){
	fmt.Println("hello world,你好，世界!")
	var a int =0
	var b float32 = 1.1
	fmt.Printf("a=%v,b=%v\n",a,b)
	fmt.Println("a=",a,",b=",b)
}


func tstLocalDecalre(){
	fmt.Println("This declare is under main function,so strong language")
}

