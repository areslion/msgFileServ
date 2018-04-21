package strategy

import(
	"testing"
	// "util"
	"log"
)

func Test_getSoftListFromDB(t *testing.T){
	var sft StrategySoft
	res,_,_ := sft.getStrategy("")
	for ix,itx:= range res.ListSft.List {
		log.Println(ix,itx.Enable,itx.Namex)
	}
}