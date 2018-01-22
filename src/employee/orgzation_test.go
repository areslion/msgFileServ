package employee
import (
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
	"testing"
	"util"
)

var manlst sxManList
func Test_readAllMan(t *testing.T){
	util.L2I("Start unit test")
	manlst.readAllMan()
	util.L2I("%d\n\n",util.SizeStruct(manlst))

	//manlst.getOrgMap()
}

func Test_getDep(t *testing.T){
	// for ix,itm :=range manlst.lstDep {
	// 	util.L2I("%d %d %v",ix,len(itm),itm)
	// }

	// for ix,itm := range manlst.lstKeyDep {
	// 	for key,valx := range itm {
	// 		util.L2I("%d %d key=%s  %v",ix,len(itm),key,valx)
	// 	}
	// }

	//for key,val := range manlst.lstKeyDep[1] {
		//util.L2I("%v %v",key,val)
		//if key == "楚雄供电局>客户服务中心" {
		//	util.L2I("%v",val)
		//}
	//}

	var orx sxOrg 
	var manx sxMan
	var manxLst sxManList

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.path = "m"+cst_sep+"a"+cst_sep+"b"+cst_sep+"c"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.path = "n"+cst_sep+"a"+cst_sep+"b"+cst_sep+"c"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.path = "n"+cst_sep+"x"+cst_sep+"b"+cst_sep+"c"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.path = "n"+cst_sep+"x"+cst_sep+"b"+cst_sep+"c1"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	log.Println("")
	log.Println("")
	log.Println("new.............")
	manx.path = "n"+cst_sep+"x"+cst_sep+"b"+cst_sep+"c2"
	manx.parse(cst_sep,true,&manxLst)
	orx.insertChild(&manx)

	orx.insertChild(&manx)
	orx.insertChild(&manx)
	orx.insertChild(&manx)

	
	bts,_:=json.Marshal(orx)
	util.L2I(string(bts))
}

func Test_orgTree(t *testing.T){
	var orgx sxOrg

	for _,itm := range manlst.mapLstMan {
		orgx.insertChild(&itm)
	}

	bts,_,bret := orgx.toJson();if bret {
		ioutil.WriteFile(".\\tree.json",bts,os.ModePerm)
	}
	
}
