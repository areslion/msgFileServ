package servMsg
import(
	"strconv"
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

const (
	cst_json = "application/json; charset=utf-8"
)

func init(){
	http.HandleFunc("/msgfile/newmsg", newmsg)            //POST upload software
	http.HandleFunc("/msgfile/usrget", usrget)            //POST upload software
	http.HandleFunc("/msgfile/usrupdate", usrupdate)      
	http.HandleFunc("/msgfile/gettsk", gettsk)
	http.HandleFunc(cst_prefix_getfil, getfile)
}

func gettsk(t_res http.ResponseWriter,t_ask *http.Request){
	util.L2I("gettsk called "+t_ask.Method)

	if t_ask.Method =="GET"{
		tskid := t_ask.FormValue("task")
		bts,_ := getOneTsk(tskid)
		t_res.Header().Set("Content-Type",cst_json)
		t_res.Write(bts)
	}
}

func getfile(t_res http.ResponseWriter,t_ask *http.Request){
	util.L2I("getfile called "+t_ask.Method)

	if t_ask.Method=="GET" {
		taskid := t_ask.FormValue("task")
		filename := t_ask.FormValue("file")

		util.L2I("%s %s",taskid,filename)

		// folder, fldID, _ := util.GetPathEle(t_ask.URL.Path)
		// prefix := "/" + folder + "/"
		// sft.FolderID = fldID
		// util.L2I("start fileservr(" + prefix + " " + sft.GetFolderPath(software.CfgSft, false) + ")")
		// staticFServ := http.StripPrefix(prefix, http.FileServer(http.Dir(sft.GetFolderPath(software.CfgSft, false))))
		// staticFServ.ServeHTTP(t_res, t_ask)
	}
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

func usrget(t_res http.ResponseWriter,t_ask *http.Request){
	util.L2I(t_ask.Method)

	if t_ask.Method == "GET" {
		idx := t_ask.FormValue("id")
		npage := t_ask.FormValue("page")
		nlimt := t_ask.FormValue("limit")
		nstatu := t_ask.FormValue("status")

		bts := getUsrMsg(idx,npage,nlimt,nstatu)
		t_res.Header().Set("Content-Type", "application/json; charset=utf-8")
		t_res.Write(bts)
	}



	// util.L2I(fmt.Sprintf("%v",t_ask.URL))

	// urlRes,_ := url.Parse(t_ask.RequestURI)
	// util.L2I(fmt.Sprintf("%v",urlRes))

	// urlV := urlRes.Query()
	// util.L2I(fmt.Sprintf("%v",urlV))
	// for ix,itm := range urlV {
	// 	util.L2I(fmt.Sprintf("%v  %v",ix,itm))
	// }

	// log.Println(t_ask.FormValue("name2"))
}

func usrupdate(t_res http.ResponseWriter,t_ask *http.Request){
	util.L2I("usrupdate "+t_ask.Method)
	if t_ask.Method == "POST" {
		tsk := t_ask.FormValue("tsk")
		dev := t_ask.FormValue("dev")
		statux,_ := strconv.Atoi(t_ask.FormValue("status"))

		updateUsrTsk(tsk,dev,statux)
	}
}

