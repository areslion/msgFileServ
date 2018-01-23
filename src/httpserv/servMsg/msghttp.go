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
var numCallTst int =0

func init(){
	http.HandleFunc("/msgfile/newmsg", newmsg)//POST upload software
	http.HandleFunc("/msgfile/usrget", usrget)//GET one usr's message task list
	http.HandleFunc("/msgfile/admget", admget)//GET administrator's message task list
	http.HandleFunc("/msgfile/usrupdate", usrupdate)//update one usr's one task status      
	http.HandleFunc("/msgfile/gettsk", gettsk)//obtain one task's detail info
	http.HandleFunc("/msgfile/tsksendlst", getsendlst)//obtain one task's detail info
	http.HandleFunc(cst_prefix_getfil, getfile)//download a file resourse

	http.HandleFunc("/msgfile/delmsg", delmsgfile)//obtain one task's detail info
}

func admget(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("admget %s",t_ask.Method)

	bret := false
	if t_ask.Method == "GET" {
		strpage := t_ask.FormValue("page")
		strlimt := t_ask.FormValue("limit")
		// npage,_ := strconv.Atoi(strpage)
		// nlimit,_ := strconv.Atoi(strlimt)
		// npage = npage*nlimit

		var bts []byte
		bts ,bret= getAdminMsg(strpage,strlimt)
		t_res.Header().Set("Content-Type", "application/json; charset=utf-8")
		t_res.Write(bts)
	}
	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}

func delmsgfile(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("delmsgfile %s %s",t_ask.Method,t_ask.URL.Path)

	bret := false
	tskid := t_ask.FormValue("task")
	if t_ask.Method=="POST"{ bret = delMsg(tskid)}
	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}

func gettsk(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("gettsk called "+t_ask.Method)

	bret := false
	if t_ask.Method =="GET"{
		tskid := t_ask.FormValue("task")

		var bts []byte
		bts,bret = getOneTskDetail(tskid)
		t_res.Header().Set("Content-Type",cst_json)
		t_res.Write(bts)
	}
	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}


func getsendlst(t_res http.ResponseWriter,t_ask *http.Request){
	numCallTst++
	util.L3I("getsendlst called %s %d",t_ask.Method,numCallTst)

	bret := false
	if t_ask.Method =="GET"{
		tskid := t_ask.FormValue("task");if len(tskid)!=36 {
			util.L4E("getsendlst invalid task guid %s",tskid)
			t_res.WriteHeader(http.StatusNotAcceptable)
			return
		}
		page ,err:= strconv.Atoi(t_ask.FormValue("page"));if err!=nil{
			util.L4E("getsendlst strconv.Atoi(page) %s",err.Error())
			t_res.WriteHeader(http.StatusNotAcceptable)
			return
		}
		limit,err := strconv.Atoi(t_ask.FormValue("limit"));if err!=nil{
			util.L4E("getsendlst strconv.Atoi(limit)",err.Error())
			t_res.WriteHeader(http.StatusNotAcceptable)
			return
		}

		var bts []byte
		bts,bret = getOneTskSendDetail(tskid,page,limit);if bret {
			t_res.Header().Set("Content-Type",cst_json)
			t_res.Write(bts)
		}
	}

	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}

func getfile(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("getfile called "+t_ask.Method)

	if t_ask.Method=="GET" {
		util.NewFileServ(t_ask,&t_res,m_cfg.ServFile.PathMsg)
	}
}

func newmsg(t_res http.ResponseWriter, t_ask *http.Request) {
	util.L3I("/msgfile/newmsg called "+t_ask.Method)
	bret := false
	if t_ask.Method =="POST" {
		bret = parseAsk(t_ask)
	}	
	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}


func parseAsk(t_ask *http.Request) (r_ret bool) {
	muti_reader, _err := t_ask.MultipartReader()
	var bfileSave, bDesc bool =true ,false 
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
	} else {
		util.L4E("parseAsk %s",_err.Error())
	}

	if bDesc {
		var msgid string
		msgid,r_ret = insertDBBytes(bufDes.Bytes());if r_ret==false {return}

		for _,itm := range attx {
			if itm.bvalid==false { break }

			bsaves := saveAttach(itm.name,msgid,itm.buf.Bytes())
			bfileSave = bfileSave||bsaves
		}
	} else {r_ret=false}

	r_ret = true
	return
}

func usrget(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	bret := false
	if t_ask.Method == "GET" {
		idx := t_ask.FormValue("id")
		npage := t_ask.FormValue("page")
		nlimt := t_ask.FormValue("limit")
		nstatu := t_ask.FormValue("status")
		tmx := t_ask.FormValue("time")

		var bts []byte
		bts ,bret= getUsrMsg(idx,npage,nlimt,nstatu,tmx)
		t_res.Header().Set("Content-Type", "application/json; charset=utf-8")
		t_res.Write(bts)
	}
	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}



	// util.L3I(fmt.Sprintf("%v",t_ask.URL))

	// urlRes,_ := url.Parse(t_ask.RequestURI)
	// util.L3I(fmt.Sprintf("%v",urlRes))

	// urlV := urlRes.Query()
	// util.L3I(fmt.Sprintf("%v",urlV))
	// for ix,itm := range urlV {
	// 	util.L3I(fmt.Sprintf("%v  %v",ix,itm))
	// }

	// log.Println(t_ask.FormValue("name2"))
}

func usrupdate(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("usrupdate "+t_ask.Method)

	bret := false
	if t_ask.Method == "POST" {
		tsk := t_ask.FormValue("task")
		dev := t_ask.FormValue("dev")
		statux,_ := strconv.Atoi(t_ask.FormValue("status"))

		util.L3I("usrupdate (task=%s dev=%s status=%d)",tsk,dev,statux)
		bret = updateUsrTsk(tsk,dev,statux)
	}
	if !bret {t_res.WriteHeader(http.StatusNotAcceptable)}
}

