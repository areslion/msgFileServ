package servMsg

import (
	"database/sql"	
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
	
	_ "github.com/go-sql-driver/mysql"
)
import (
	"os"
	"dbbase"
	"util"
)


type sxReciever struct {
	Guid string `json:"guid"`
}
type sxAttatche struct {
	Index int    `json:"index"`
	Name  string `json:"filename"`
	Url   string `json:"url"`
}
type sxExctm struct {
	Tmx string `json:"tmx"`
	Tmy string `json:"tmy"`
}
type sxMsg struct {
	Name     string       `json:"name"`
	Guid     string       `json:"guid"`
	Tmx      string       `json:"tmx"`
	Tmy      string       `json:"tmy"`
	Sender   string       `json:"sender"`
	Os       int          `json:"os"`
	Auto     int          `json:"auto"`
	Popup    int          `json:"popupwindow"`
	Desc     string       `json:"desc"`
	Reciever []sxReciever `json:"reciever"`
	Attach   []sxAttatche `json:"attachement"`
	Exctm    []sxExctm    `json:"tmExc"`
}

var m_cfg *util.SxCfgAll



func StartServMsg(t_cfg *util.SxCfgAll) {
	m_cfg = t_cfg

	Tstservmsg()
}


func insertDBBytes(t_bts []byte)(r_id string,r_ret bool){
	var msgx *sxMsg
	msgx,r_ret = parseJson(t_bts)	
	if r_ret == false { 
		util.L3E("insertDBBytes fail to parseJson")
		return
	 }

	r_id = msgx.Guid
	r_ret = insertDB(msgx)
	
	 return
}


func insertDB(t_msg *sxMsg)(r_ret bool){
	var cnt *sql.DB
	cnt, r_ret = dbbase.Open(m_cfg)
	if r_ret == false {
		return
	}
	defer dbbase.Close()
	defer cnt.Close()

	r_ret = insertAbstract(cnt,t_msg);if r_ret==false {return}
	r_ret = insertSender(cnt,t_msg);if r_ret==false {return}
	r_ret = saveDescMsg(t_msg);if r_ret==false {return}

	return
}
//Insert message into database
func insertAbstract(t_cnn *sql.DB,t_msg *sxMsg) (b_ret bool) {
	sqlcmd := "REPLACE INTO msgAbstract (numMsg,namex,tmx,tmy,tmm,os,autoexe,popup,numSender)"
	sqlcmd += "VALUE(?,?,?,?,?,?,?,?,?)"
	smt, err := t_cnn.Prepare(sqlcmd)
	if err != nil {
		util.L3E("insertAbstract fail to Prepare " + err.Error())
		return
	}

	var px *sxMsg = t_msg
	tmNow := time.Now().Format("2006-01-02 15:04:05")
	_, err = smt.Exec(px.Guid, px.Name, px.Tmx, px.Tmy, tmNow, px.Os, px.Auto, px.Popup, px.Sender)
	if err != nil {
		util.L4F("insertAbstract " + err.Error())
		return
	}

	b_ret = true
	util.L2I("insertAbstract add a msg task " + t_msg.Guid + " " + t_msg.Name)
	return
}

func insertSender(t_cnn *sql.DB,t_msg *sxMsg) (b_ret bool) {
	sqlcmd := "REPLACE INTO msgSend (numMsg,numReciever,tmm)"
	sqlcmd += "VALUE(?,?,?)"
	smt, err := t_cnn.Prepare(sqlcmd)
	if err != nil {
		util.L3E("insertSender  fail to Prepare " + err.Error())
		return
	}

	var px *sxMsg = t_msg
	tmNow := time.Now().Format("2006-01-02 15:04:05")

	for ix,item := range t_msg.Reciever{
		_, err = smt.Exec(px.Guid,item.Guid,tmNow)
		if err != nil {
			util.L4F("saveNewSft  " + err.Error())
			return
		} else{
			util.L1T("add a task to a user "+ fmt.Sprintf("%v  %d",item,ix))
		}
	}

	b_ret = true
	util.L2I("insertSender add a msg " + t_msg.Guid + " " + t_msg.Name+ " to usr="+fmt.Sprintf("%d",len(px.Reciever)))
	return
}

func parseJson(t_bts []byte) (r_msg *sxMsg,b_ret bool) {
	r_msg = new(sxMsg)
	err := json.Unmarshal(t_bts, r_msg)
	if err != nil {
		util.L3E("parseJson  " + err.Error())		
		return
	} else {b_ret = true}
	util.L1T(fmt.Sprintf("%v", r_msg))

	return
}

func Tstservmsg() {
	bts, err := ioutil.ReadFile(".\\cfg\\tsms_msg.json")
	if err == nil {
		//log.Println(string(bts))
		msgx,_ := parseJson(bts)
		insertDB(msgx)
	}
}

func saveAttach(t_name string,t_id string,t_bts []byte)(b_ret bool){
	folder := m_cfg.ServFile.PathMsg + t_id + util.GetOSSeptor()
	strAtc := fmt.Sprintf("%v",t_name)

	os.MkdirAll(folder,0711)
	_,b_ret = util.SaveFileBytes(folder+t_name,t_bts)
	if b_ret==true {  util.L1T("saveAttach OK "+strAtc) } else {
		util.L1T("saveAttach KO "+strAtc)
	}

	return
}

func saveDescMsg(t_msg *sxMsg)(r_ret bool){
	folder := m_cfg.ServFile.PathMsg + t_msg.Guid + util.GetOSSeptor()
	os.MkdirAll(folder,0711)

	bts,_ := json.Marshal(t_msg)
	_,r_ret = util.SaveFileBytes(folder+"Desc.json",bts)
	if r_ret==true {  util.L1T("saveDesc OK ") } else {
		util.L1T("saveDesc KO ")
	}

	return
}

func saveDescBytes(t_id string,t_bts []byte)(r_ret bool){
	folder := m_cfg.ServFile.PathMsg + t_id + util.GetOSSeptor()
	os.MkdirAll(folder,0711)

	_,r_ret = util.SaveFileBytes(folder,t_bts)
	if r_ret==true {  util.L1T("saveDesc OK ") } else {
		util.L1T("saveDesc KO ")
	}

	return
}