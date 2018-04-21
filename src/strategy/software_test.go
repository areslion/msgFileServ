package strategy

import(
	"testing"
	// "util"
	"log"
)

func Test_getSoftListFromDB(t *testing.T){
	var sft StrategySoft
	list := sft.getSoftListFromDB()
	for ix,itx:= range list.List {
		log.Println(ix,itx.Enable,itx.Namex)
	}
}