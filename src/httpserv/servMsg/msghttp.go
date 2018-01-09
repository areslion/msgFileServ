package servMsg
import(
	"fmt"
	"io"
	"net/http"
)
import (
	"util"
)


func init(){
	http.HandleFunc("/msgfile/newmsg", newmsg)            //POST upload software
}

func newmsg(t_res http.ResponseWriter, t_ask *http.Request) {
	util.L2I("/msgfile/newmsg called")
	parseAsk(t_ask)
}


func parseAsk(t_ask *http.Request) bool {
	muti_reader, _err := t_ask.MultipartReader()
	var bfileSave bool = false

	if _err == nil {
		for {
			part, err := muti_reader.NextPart()
			if err == io.EOF {
				break
			}
			
			util.L1T(fmt.Sprintf("%v",part))
		}
	}

	return bfileSave
}