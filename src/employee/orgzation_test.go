package employee
import (
	"log"
	"testing"
	"util"
)

var manlst sxManList
func Test_readAllMan(t *testing.T){
	util.L3I("Start unit test")
	manlst.readAllMan()
	util.L3I("%d\n\n",util.SizeStruct(manlst))
}

func Test_getDep(t *testing.T){
	var orx sxOrg 
	var manx sxMan
	var manxLst sxManList

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "m"+getSep()+"a"+getSep()+"b"+getSep()+"c"
	manx.parse(getSep(),true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "n"+getSep()+"a"+getSep()+"b"+getSep()+"c"
	manx.parse(getSep(),true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "n"+getSep()+"x"+getSep()+"b"+getSep()+"c"
	manx.parse(getSep(),true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "n"+getSep()+"x"+getSep()+"b"+getSep()+"c1"
	manx.parse(getSep(),true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Name = "C2.0"
	manx.Path = "n"+getSep()+"x"+getSep()+"b"+getSep()+"c2"
	manx.parse(getSep(),true,&manxLst)
	orx.insertChild(&manx)

	manx.Name = "C2.1"
	orx.insertChild(&manx)
	manx.Name = "C2.2"
	orx.insertChild(&manx)
	manx.Name = "C2.3"	
	orx.insertChild(&manx)

	
	bts,_,_:=orx.toJson()
	util.L3I(string(bts))


	lst := orx.GetLstDepat("n"+getSep()+"x"+getSep()+"b")
	for ix,itm := range lst{
		util.L3I("%d %s",ix,itm)
	}
}

var orgx sxOrg
func Test_orgTree(t *testing.T){
	for _,itm := range manlst.lstMan {
		orgx.insertChild(&itm)
	}

	orgx.saveJson(".\\tree.json")
	
}


func Test_GetMsg(t *testing.T){
	lst,_ := orgx.GetLstDepat("")
	for ix,itm := range lst{
		util.L3I("%d %s",ix,itm)
	}

	//var strIn string
	//fmt.Println("输入部门名称:")
	//fmt.Scanln(strIn)
	lst,_ = orgx.GetLstDepat("楚雄供电局")
	//lst = orgx.GetLstDepat(strIn)
	for ix,itm := range lst{
		util.L3I("%d %s",ix,itm)
	}
}