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


func insertDB(t_msg *sxMsg){
	cnt, bret := dbbase.Open(m_cfg)
	if bret == false {
		return
	}
	defer dbbase.Close()
	defer cnt.Close()

	insertAbstract(cnt,t_msg)
	insertSender(cnt,t_msg)
}
//Insert message into database
func insertAbstract(t_cnn *sql.DB,t_msg *sxMsg) (b_ret bool) {
	sqlcmd := "REPLACE INTO msgAbstract (numMsg,namex,tmx,tmy,tmm,os,autoexe,popup,numSender)"
	sqlcmd += "VALUE(?,?,?,?,?,?,?,?,?)"
	smt, err := t_cnn.Prepare(sqlcmd)
	if err != nil {
		util.L3E("InsertDB  fail to Prepare " + err.Error())
		return
	}

	var px *sxMsg = t_msg
	tmNow := time.Now().Format("2006-01-02 15:04:05")
	_, err = smt.Exec(px.Guid, px.Name, px.Tmx, px.Tmy, tmNow, px.Os, px.Auto, px.Popup, px.Sender)
	if err != nil {
		util.L4F("saveNewSft  " + err.Error())
		return
	}

	b_ret = true
	util.L2I("saveNewSft  add a msg task " + t_msg.Guid + " " + t_msg.Name)
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

func parseJson(t_bts []byte) (r_msg *sxMsg) {
	r_msg = new(sxMsg)
	err := json.Unmarshal(t_bts, r_msg)
	if err != nil {
		util.L3E("parseJson  " + err.Error())
		return
	}
	util.L1T(fmt.Sprintf("%v", r_msg))

	return
}

func Tstservmsg() {
	bts, err := ioutil.ReadFile(".\\cfg\\tsms_msg.json")
	if err == nil {
		//log.Println(string(bts))
		msgx := parseJson(bts)
		insertDB(msgx)
	}
}

func saveFileBytes( buf []byte) (r_path string, b_ret bool) {
	return 
}
