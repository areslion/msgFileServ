package employee
import (
	"io/ioutil"
	"log"
	_ "encoding/json"
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
