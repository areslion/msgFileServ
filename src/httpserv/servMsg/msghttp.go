package servMsg
import(
	"bytes"
	"fmt"
	"io"
	"net/http"
)
import (
	"util"
)

type sxAttEle struct{
	name string
	bvalid bool
	buf *bytes.Buffer
}


func init(){
	http.HandleFunc("/msgfile/newmsg", newmsg)            //POST upload software
}

func newmsg(t_res http.ResponseWriter, t_ask *http.Request) {
	util.L2I("/msgfile/newmsg called "+t_ask.Method)
	if t_ask.Method =="POST" {
		parseAsk(t_ask)
	}	
}


func parseAsk(t_ask *http.Request) (r_ret bool) {
	muti_reader, _err := t_ask.MultipartReader()
	var bfileSave, bDesc bool =false ,false 
	var attx [100]sxAttEle
	var ixat int =0
	var bufDes = new(bytes.Buffer)

	if _err == nil {
		for {
			part, err := muti_reader.NextPart()
			if err == io.EOF {
				break
			}
			util.L1T(fmt.Sprintf("%v",part))

			if part.FormName() =="attachment"{
				if ixat>=100 { continue }

				attx[ixat].buf = new(bytes.Buffer)
				attx[ixat].buf.ReadFrom(part)
				attx[ixat].name = part.FileName()
				attx[ixat].bvalid = true
				ixat++
			} else if part.FormName() =="description" {
				bufDes.ReadFrom(part)
				bDesc = true
			}


		}
	}

	if bDesc && ixat >0 {
		var msgid string
		msgid,r_ret = insertDBBytes(bufDes.Bytes());if r_ret==false {return}

		for _,itm := range attx {
			if itm.bvalid==false { break }

			saveAttach(itm.name,msgid,itm.buf.Bytes())
		}
	}

	return bfileSave
}