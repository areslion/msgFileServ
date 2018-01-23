package employee
import (
	"log"
	"testing"
	"util"
)

var manlst sxManList
func Test_readAllMan(t *testing.T){
	util.L2I("Start unit test")
	manlst.readAllMan()
	util.L2I("%d\n\n",util.SizeStruct(manlst))
}

func Test_getDep(t *testing.T){
	var orx sxOrg 
	var manx sxMan
	var manxLst sxManList

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "m"+cst_sep+"a"+cst_sep+"b"+cst_sep+"c"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "n"+cst_sep+"a"+cst_sep+"b"+cst_sep+"c"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "n"+cst_sep+"x"+cst_sep+"b"+cst_sep+"c"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Path = "n"+cst_sep+"x"+cst_sep+"b"+cst_sep+"c1"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.Name = "C2.0"
	manx.Path = "n"+cst_sep+"x"+cst_sep+"b"+cst_sep+"c2"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	manx.Name = "C2.1"
	orx.insertChild(&manx)
	manx.Name = "C2.2"
	orx.insertChild(&manx)
	manx.Name = "C2.3"	
	orx.insertChild(&manx)

	
	bts,_,_:=orx.toJson()
	util.L2I(string(bts))


	lst := orx.GetLstDepat("n"+cst_sep+"x"+cst_sep+"b")
	for ix,itm := range lst{
		util.L2I("%d %s",ix,itm)
	}
}

var orgx sxOrg
func Test_orgTree(t *testing.T){
	for _,itm := range manlst.mapLstMan {
		orgx.insertChild(&itm)
	}

	orgx.saveJson(".\\tree.json")
	
}


func Test_GetMsg(t *testing.T){
	lst := orgx.GetLstDepat("")
	for ix,itm := range lst{
		util.L2I("%d %s",ix,itm)
	}

	lst = orgx.GetLstDepat("楚雄供电局>物流服务中心")
	for ix,itm := range lst{
		util.L2I("%d %s",ix,itm)
	}
}