package servMsg

import (
	"strconv"
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


const (
	cst_tsks_u = 0x0001 //初始未执行状态
	cst_tsks_r = 0x0002 //客户端收到任务消息
	cst_tsks_e = 0x0004 //客户端执行成功
	cst_tsks_f = 0x0008 //客户端执行失败
)

type sxReciever struct {
	Guid string `json:"guid"`
}
type sxAttatche struct {
	Index int    `json:"index"`
	Name  string `json:"filename"`
	Url   string `json:"url"`
	Descx string `json:"desc"`
	Sizex string `json:"size"`
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
	Tmexc    string       `json:"tmexc"`
	Sender   string       `json:"sender"`
	Os       int          `json:"os"`
	Auto     int          `json:"auto"`
	Popup    int          `json:"popupwindow"`
	Desc     string       `json:"desc"`
	Status	 int		  `json:"status"`
	Reciever []sxReciever `json:"reciever"`
	Attach   []sxAttatche `json:"attachement"`
	Exctm    []sxExctm    `json:"tmExc"`
}
type sxMsgAskRes struct {
	Totalnum int `json:"taotlnum"`
	Page int `json:"page"`
	Limit int `json:"limit"`
	UsrID string `json:"usrid"`
	Lst []sxMsg
}

var m_cfg *util.SxCfgAll

func closex(t_cnn *sql.DB){
	if t_cnn!=nil {t_cnn.Close()}
	dbbase.Close()
}

func getUsrMsg(t_id,t_page,t_limit,t_status string)(r_bts []byte){
	cnn := openx();if cnn==nil {return}
	defer closex(cnn)
	
	nstatus,err := strconv.Atoi(t_status); if err != nil {
		util.L3E("getUsrMsg strconv.Atoi(t_status) "+err.Error())
		return
	}

	var strFlag string
	var num int
	fx := func (x_val int){
		if (nstatus & x_val)>0 {
			if num>0 {strFlag += " OR "}
			strFlag += "statusx="+fmt.Sprintf("%d ",x_val) 
			num++
		}
	}
	fx(cst_tsks_e)
	fx(cst_tsks_f)
	fx(cst_tsks_r)
	fx(cst_tsks_u)
	if len(strFlag)>0 { strFlag = "numReciever="+"'"+t_id+"'"+" AND "+ "("+strFlag+") " } else {
		strFlag = "numReciever="+"'"+t_id+"'"
	}
	

	sqlcmd := "SELECT (SELECT COUNT(*) FROM msgSend WHERE "+strFlag+")num,namex,statusx,tmx,tmy,tmExc,descx "
	sqlcmd += "FROM msgSend "
	sqlcmd += "WHERE "+strFlag
	sqlcmd += "LIMIT ?,? "

	util.L2I(sqlcmd)
	smt,err := cnn.Prepare(sqlcmd);if err!=nil {
		util.L3E("getUsrMsg Prepare "+err.Error())
	}
	rows,err := smt.Query(t_page,t_limit); if err !=nil {
		util.L3E("getUsrMsg smt.Query "+err.Error())
		return
	}


	var resMsg sxMsgAskRes
	resMsg.Page,_ = strconv.Atoi(t_page)
	resMsg.Limit,_ = strconv.Atoi(t_limit)
	resMsg.UsrID = t_id

	for rows.Next() {
		var ele sxMsg
		rows.Scan(&resMsg.Totalnum,&ele.Name,&ele.Status,&ele.Tmx,&ele.Tmy,&ele.Tmexc,&ele.Desc)
		resMsg.Lst = append(resMsg.Lst,ele)

		util.L2I(fmt.Sprintf("%v",ele))
		//if rows.Next()==false {break}
		//util.L2I(fmt.Sprintf("%v",resMsg.Lst))
	}
	bts,err := json.Marshal(resMsg); if err !=nil {
		util.L3E("getUsrMsg json.Marshal(resMsg)"+err.Error())
		return
	}
	r_bts = bts

	util.L2I("getUsrMsg "+t_id+"("+fmt.Sprintf("toatl=%d  %d/(Limit %d %d) %d",resMsg.Totalnum,resMsg.Page,resMsg.Page,resMsg.Limit,len(resMsg.Lst))+")")

	return
}

func openx()(r_cnnt *sql.DB){
	var bret bool
	r_cnnt, bret = dbbase.Open(m_cfg)
	if bret == false {
		r_cnnt = nil
		return
	}
	return
}

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