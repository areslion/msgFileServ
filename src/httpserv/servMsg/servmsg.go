package servMsg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)
import (
	"dbbase"
	"util"
)

const (
	cst_tsks_1u = 0x0001 //初始未执行状态
	cst_tsks_2r = 0x0002 //客户端收到任务消息
	cst_tsks_3e = 0x0004 //客户端执行成功
	cst_tsks_4f = 0x0008 //客户端执行失败

	cst_fix_desc      = "Desc.json"
	cst_prefix_getfil = "/msgfile/getfile/"
)

var cst_tsksArr = [...]int{cst_tsks_1u, cst_tsks_2r, cst_tsks_3e, cst_tsks_4f}



type (
	sxReciever struct {
		Guid string `json:"guid"`
	}
	sxAttatche struct {
		Index int    `json:"index"`
		Name  string `json:"filename"`
		Url   string `json:"url"`
		Descx string `json:"desc"`
		Sizex string `json:"size"`
	}
	sxExctm struct {
		Tmx string `json:"tmx"`
		Tmy string `json:"tmy"`
	}
	sxMsg struct {
		Name     string       `json:"name"`
		Guid     string       `json:"guid"`
		Tmx      string       `json:"tmx"`
		Tmy      string       `json:"tmy"`
		Tmexc    string       `json:"tmexcok"`
		Sender   string       `json:"sender"`
		Os       int          `json:"os"`
		Auto     int          `json:"auto"`
		Popup    int          `json:"popupwindow"`
		Desc     string       `json:"desc"`
		Status   int          `json:"status"`
		Numsend  int          `json:"numSend"`
		NumOK    int          `json:"numOK"`
		NumKO    int          `json:"numKO"`
		Reciever []sxReciever `json:"reciever"`
		Attach   []sxAttatche `json:"attachement"`
		Exctm    []sxExctm    `json:"tmExc"`
	}
	sxMsgAskRes struct {
		Totalnum int    `json:"taotlnum"`
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		UsrID    string `json:"usrid"`
		Lst      []sxMsg
	}

	sxOneReciever struct{
		NumDev string `json:"uuid"`
		Status int `json:"status"`
		TmExc string `json:"tmExc"`
		Os int `json:"os"`
		OwnerName string `json:"ownername"`
		OwnerDepart string `json:"ownerDepart"`
		Detail string `json:"detail"`
	}
	sxTskSendDetial struct {
		Totalnum int `json:"taotlnum"`
		Page int `json:"page"`
		Limit int `json:"limit"`
		Task string `json:"guid"`
		Name string `json:"name"`
		Detail string `json:"detail"`
		NAll int `json:"numAll"`
		NOK int `json:"numOK"`
		NKO int `json:"numKO"`
		Lst []sxOneReciever `json:"list"`
	}
)

var m_cfg *util.SxCfgAll

func closex(t_cnn *sql.DB) {
	if t_cnn != nil {
		t_cnn.Close()
	}
	dbbase.Close(t_cnn)
}

func delMsg(t_msgid string) (b_ret bool) {
	if len(t_msgid) != 36 {
		util.L3E("delMsg invalid msgid(%s=%d)", t_msgid, len(t_msgid))
		return
	}

	dbopt,bret := dbbase.NewSxDB(&m_cfg.Db,"delMsg");if !bret {return}
	defer dbopt.Close()

	util.L2I("delMsg start to delete msg " + t_msgid)
	folder := m_cfg.ServFile.PathMsg + util.GetOSSeptor() + t_msgid
	if util.IsExists(folder) {
		err := os.RemoveAll(folder)
		if err != nil {
			util.L3E("delMsg os.RemoveAll(%s) %s", folder, err.Error())
			return
		} else {
			util.L2I("delMsg foler %s removed", folder)
		}
	}

	dbopt.Sqlcmd = "DELETE FROM msgAbstract WHERE numMsg=?"
	if !dbopt.Exc(t_msgid){return}
	util.L2I("delMsg deleted %s form msgAbstract %d",t_msgid,dbopt.Affected())

	dbopt.Sqlcmd = "DELETE FROM msgSend WHERE numMsg=?"
	if !dbopt.Exc(t_msgid) {return}
	util.L2I("delMsg deleted %s form msgSend %d",t_msgid,dbopt.Affected())

	b_ret = true
	return
}

func getAdminMsg(t_page, t_limit string) (r_bts []byte, b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&m_cfg.Db,"getAdminMsg") ; if !bret {return}
	defer dbopt.Close()

	npage ,_:= strconv.Atoi(t_page)
	nlimit ,_:= strconv.Atoi(t_limit)

	dbopt.Sqlcmd = "SELECT (SELECT COUNT(*) FROM msgAbstract)num,namex,tmx,tmy,numSent,numSentOK,numSentKO,numMsg,os "
	dbopt.Sqlcmd  += "FROM msgAbstract "
	dbopt.Sqlcmd  += "LIMIT ?,? "

	if !dbopt.Query(npage*nlimit, t_limit){return}

	var resMsg sxMsgAskRes
	resMsg.Page = npage
	resMsg.Limit = nlimit

	for dbopt.Rows.Next() {
		var ele sxMsg
		dbopt.Rows.Scan(&resMsg.Totalnum, &ele.Name, &ele.Tmx, &ele.Tmy, &ele.Numsend, &ele.NumOK, &ele.NumKO,&ele.Guid,&ele.Os)
		resMsg.Lst = append(resMsg.Lst, ele)

		util.L2I(fmt.Sprintf("%v", ele))
	}
	bts, err := json.Marshal(resMsg)
	if err != nil {
		util.L3E("getAdminMsg json.Marshal(resMsg)" + err.Error())
		return
	}
	r_bts = bts

	b_ret = true
	util.L2I("getAdminMsg " + "(" + fmt.Sprintf("toatl=%d  %d/(Limit %d %d) %d", resMsg.Totalnum, resMsg.Page, resMsg.Page, resMsg.Limit, len(resMsg.Lst)) + ")")

	return
}

func getUsrMsg(t_id, t_page, t_limit, t_status string) (r_bts []byte, b_ret bool) {
	dbopt,bret:=dbbase.NewSxDB(&m_cfg.Db,"getUsrMsg") ; if !bret {return}
	defer dbopt.Close()

	nstatus, _:= strconv.Atoi(t_status)
	npage,_:=strconv.Atoi(t_page)
	nlimit,_:=strconv.Atoi(t_limit)

	var strFlag string
	var num int
	fx := func(x_val int) {
		if (nstatus & x_val) > 0 {
			if num > 0 {
				strFlag += " OR "
			}
			strFlag += "statusx=" + fmt.Sprintf("%d ", x_val)
			num++
		}
	}
	fx(cst_tsks_3e)
	fx(cst_tsks_4f)
	fx(cst_tsks_2r)
	fx(cst_tsks_1u)
	if len(strFlag) > 0 {
		strFlag = "numReciever=" + "'" + t_id + "'" + " AND " + "(" + strFlag + ") "
	} else {
		strFlag = "numReciever=" + "'" + t_id + "'"
	}

	dbopt.Sqlcmd = "SELECT (SELECT COUNT(*) FROM msgSend WHERE " + strFlag + ")num,namex,statusx,tmx,tmy,tmExc,descx "
	dbopt.Sqlcmd += "FROM msgSend "
	dbopt.Sqlcmd += "WHERE " + strFlag
	dbopt.Sqlcmd += "LIMIT ?,? "

	util.L1T(dbopt.Sqlcmd)
	if !dbopt.Query(npage*nlimit, t_limit){return}

	var resMsg sxMsgAskRes
	resMsg.Page = npage
	resMsg.Limit = nlimit
	resMsg.UsrID = t_id

	for dbopt.Rows.Next() {
		var ele sxMsg
		dbopt.Rows.Scan(&resMsg.Totalnum, &ele.Name, &ele.Status, &ele.Tmx, &ele.Tmy, &ele.Tmexc, &ele.Desc)
		resMsg.Lst = append(resMsg.Lst, ele)

		util.L1T(fmt.Sprintf("%v", ele))
	}
	bts, err := json.Marshal(resMsg)
	if err != nil {
		util.L3E("getUsrMsg json.Marshal(resMsg)" + err.Error())
		return
	}
	r_bts = bts

	b_ret = true
	util.L2I("getUsrMsg %s(toatl=%d  %d/(Limit %d %d) %d",t_id,resMsg.Totalnum, resMsg.Page, resMsg.Page, resMsg.Limit, len(resMsg.Lst))

	return
}

func getOneTskDetail(t_msgid string) (r_bts []byte, b_ret bool) {
	path := m_cfg.ServFile.PathMsg + util.GetOSSeptor() + t_msgid + util.GetOSSeptor() + cst_fix_desc

	bts, err := ioutil.ReadFile(path)
	if err != nil {
		util.L3E(fmt.Sprintf("getOneTsk ioutil.ReadFile %s failed %s", path, err.Error()))
		return
	}

	var msgx sxMsg
	err = json.Unmarshal(bts, &msgx)
	if err != nil {
		util.L3E(fmt.Sprintf("getOneTsk json.Unmarshal(bts,&sxMsg) %s %s", path, err.Error()))
		return
	}

	b_ret = true

	util.L1T("%v", msgx)
	r_bts = []byte(fmt.Sprintf("%v", msgx))
	util.L2I("getOneTsk " + t_msgid + " OK")
	return
}


func getOneTskSendDetail(t_tsk string,t_page,t_limit int) (r_bts []byte, b_ret bool) {
	var (
		sdx sxTskSendDetial
		err error
	)

	dbopt ,bret:= dbbase.NewSxDB(&m_cfg.Db,"getOneTskSendDetail");if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd =  "SELECT namex,numSent,numSentOK,numSentKO,descx FROM msgAbstract WHERE numMsg=? " 
	if !dbopt.Query(t_tsk) {return}
	if dbopt.Rows.Next() {dbopt.Rows.Scan(&sdx.Name,&sdx.NAll,&sdx.NOK,&sdx.NKO,&sdx.Detail)} else {
		util.L3E("getOneTskSendDetail fail to get task %s",t_tsk)
		return
	}

	dbopt.Sqlcmd = "SELECT (SELECT COUNT(*) FROM msgSend WHERE numMsg=?)num,numReciever,statusx,tmExc FROM msgSend WHERE numMsg=? limit ?,?" 
	if !dbopt.Query(t_tsk,t_tsk,t_page*t_limit,t_limit){return}
	for dbopt.Rows.Next(){
		var ele sxOneReciever
		dbopt.Rows.Scan(&sdx.Totalnum,&ele.NumDev,&ele.Status,&ele.TmExc)
		sdx.Lst = append(sdx.Lst,ele)
	}

	r_bts,err = json.Marshal(&sdx);if err!=nil {
		util.L3E("getOneTskSendDetail json.Marshal(&sdx) %s",err.Error())
		return
	}

	b_ret = true
	util.L2I("getOneTskSendDetail (page%d/%d res=%d/%d)",t_page,t_limit,sdx.NAll,sdx.Totalnum)
	return
}

func openx() (r_cnnt *sql.DB) {
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

func insertDBBytes(t_bts []byte) (r_id string, r_ret bool) {
	var msgx *sxMsg
	msgx, r_ret = parseJson(t_bts)
	if r_ret == false {
		util.L3E("insertDBBytes fail to parseJson")
		return
	}

	r_id = msgx.Guid

	for ix, _ := range msgx.Attach {
		msgx.Attach[ix].Url = cst_prefix_getfil + "/" + msgx.Guid + "/" + msgx.Attach[ix].Name
		util.L2I("insertDBBytes Url=" + msgx.Attach[ix].Url)
	}
	r_ret = insertDB(msgx)

	return
}

func insertDB(t_msg *sxMsg) (r_ret bool) {
	var cnt *sql.DB
	cnt, r_ret = dbbase.Open(m_cfg)
	if r_ret == false {
		return
	}
	defer dbbase.Close(cnt)

	r_ret = insertAbstract(cnt, t_msg)
	if r_ret == false {
		return
	}
	r_ret = insertSender(cnt, t_msg)
	if r_ret == false {
		return
	}
	r_ret = saveDescMsg(t_msg)
	if r_ret == false {
		return
	}

	return
}

//Insert message into database
func insertAbstract(t_cnn *sql.DB, t_msg *sxMsg) (b_ret bool) {
	sqlcmd := "REPLACE INTO msgAbstract (numMsg,namex,tmx,tmy,tmm,os,autoexe,popup,numSender,numSent,descx)"
	sqlcmd += "VALUE(?,?,?,?,?,?,?,?,?,?,?)"
	smt, err := t_cnn.Prepare(sqlcmd)
	if err != nil {
		util.L3E("insertAbstract fail to Prepare " + err.Error())
		return
	}

	var px *sxMsg = t_msg
	tmNow := time.Now().Format("2006-01-02 15:04:05")
	_, err = smt.Exec(px.Guid, px.Name, px.Tmx, px.Tmy, tmNow, px.Os, px.Auto, px.Popup, px.Sender, len(px.Reciever), px.Desc)
	if err != nil {
		util.L4F("insertAbstract " + err.Error())
		return
	}

	b_ret = true
	util.L2I("insertAbstract add a msg task " + t_msg.Guid + " " + t_msg.Name)
	return
}

func insertSender(t_cnn *sql.DB, t_msg *sxMsg) (b_ret bool) {
	sqlcmd := "REPLACE INTO msgSend (numMsg,numReciever,tmm,tmx,tmy,namex,descx)"
	sqlcmd += "VALUE(?,?,?,?,?,?,?)"
	smt, err := t_cnn.Prepare(sqlcmd)
	if err != nil {
		util.L3E("insertSender  fail to Prepare " + err.Error())
		return
	}

	var px *sxMsg = t_msg
	tmNow := time.Now().Format("2006-01-02 15:04:05")

	for ix, item := range t_msg.Reciever {
		if len(item.Guid) != 32 {
			util.L3E("insertSender invalid receivere=" + item.Guid)
			continue
		}
		_, err = smt.Exec(px.Guid, item.Guid, tmNow, px.Tmx, px.Tmy, px.Name, px.Desc)
		if err != nil {
			util.L4F("saveNewSft  " + err.Error())
			return
		} else {
			util.L1T("add a task to a user " + fmt.Sprintf("%v  %d", item, ix))
		}
	}

	b_ret = true
	util.L2I("insertSender add a msg " + t_msg.Guid + " " + t_msg.Name + " to usr=" + fmt.Sprintf("%d", len(px.Reciever)))
	return
}

func parseJson(t_bts []byte) (r_msg *sxMsg, b_ret bool) {
	r_msg = new(sxMsg)
	err := json.Unmarshal(t_bts, r_msg)
	if err != nil {
		util.L3E("parseJson  " + err.Error())
		return
	} else {
		b_ret = true
	}
	util.L1T(fmt.Sprintf("%v", r_msg))

	return
}

func Tstservmsg() {
	bts, err := ioutil.ReadFile(".\\cfg\\tsms_msg.json")
	if err == nil {
		//log.Println(string(bts))
		msgx, _ := parseJson(bts)
		insertDB(msgx)
	}
}

func saveAttach(t_name string, t_id string, t_bts []byte) (b_ret bool) {
	folder := m_cfg.ServFile.PathMsg + t_id + util.GetOSSeptor()
	strAtc := fmt.Sprintf("%v", t_name)

	os.MkdirAll(folder, 0711)
	_, b_ret = util.SaveFileBytes(folder+t_name, t_bts)
	if b_ret == true {
		util.L1T("saveAttach OK " + strAtc)
	} else {
		util.L1T("saveAttach KO " + strAtc)
	}

	return
}

func saveDescMsg(t_msg *sxMsg) (r_ret bool) {
	folder := m_cfg.ServFile.PathMsg + t_msg.Guid + util.GetOSSeptor()
	os.MkdirAll(folder, 0711)

	util.L1T(fmt.Sprintf("%v", t_msg))
	bts, _ := json.Marshal(t_msg)
	_, r_ret = util.SaveFileBytes(folder+cst_fix_desc, bts)
	if r_ret == true {
		util.L1T("saveDesc OK ")
	} else {
		util.L1T("saveDesc KO ")
	}

	return
}

func saveDescBytes(t_id string, t_bts []byte) (r_ret bool) {
	folder := m_cfg.ServFile.PathMsg + t_id + util.GetOSSeptor()
	os.MkdirAll(folder, 0711)

	_, r_ret = util.SaveFileBytes(folder, t_bts)
	if r_ret == true {
		util.L1T("saveDesc OK ")
	} else {
		util.L1T("saveDesc KO ")
	}

	return
}

func updateUsrTsk(t_tskid, t_usrid string, t_status int) (b_ret bool) {
	var nOK, nKO int = 0, 0
	nOK, nKO, b_ret = updateSendStatus(t_tskid, t_usrid, t_status)
	if !b_ret {	return}
	if t_status >= cst_tsks_3e {
		b_ret = updateAbstractNum(t_tskid, nOK, nKO)
	}

	return
}

func updateSendStatus(t_tskid, t_usrid string, t_status int) (r_OK, r_KO int, b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&m_cfg.Db,"updateSendStatus");if !bret {return}
	defer dbopt.Close()

	bvalid := false
	for _, itm := range cst_tsksArr {
		if itm == t_status {
			bvalid = true
			break
		}
	}
	if bvalid == false {
		util.L3E("updateSendStatus invalid status (%s %s %d)",t_tskid, t_usrid, t_status)
		return
	}

	var statusOld int

	dbopt.Sqlcmd = "SELECT statusx FROM msgSend WHERE numMsg=? AND numReciever=? "
	if !dbopt.Query(t_tskid, t_usrid) {return}
	if dbopt.Rows.Next() { dbopt.Rows.Scan(&statusOld) } else {
		util.L3E("updateSendStatus smt.Query(%s) failed", dbopt.Sqlcmd)
		return
	}

	if t_status <= cst_tsks_1u {
		util.L2I("updateSendStatus invalid status(%s.%s %d)", t_usrid, t_tskid, statusOld)
		return
	}
	if statusOld == cst_tsks_3e {
		util.L2I("updateSendStatus has been set as(%s.%s %d)", t_usrid, t_tskid, statusOld)
		return
	}

	if statusOld == cst_tsks_4f && t_status == cst_tsks_3e {
		r_OK = 1
		r_KO = -1
	} else if t_status == cst_tsks_4f {
		r_OK = 0
		r_KO = 1
	} else {
		r_OK = 1
	}

	dbopt.Sqlcmd = "UPDATE msgSend SET statusx=? WHERE numMsg=? AND numReciever=? "
	if !dbopt.Exc(t_status, t_tskid, t_usrid) {return}

	affcted := dbopt.Affected()
	if affcted < 1 {
		util.L3E(fmt.Sprintf("updateUsrTsk RowsAffected()=%d no relative record(dev=%s tsk=%s) find", affcted, t_usrid, t_tskid))
		return
	}

	b_ret = true
	util.L2I("updateSendStatus OK dev=%s tsk=%s status=%d affcted=%d",t_usrid,t_tskid, t_status, affcted)

	return
}

func updateAbstractNum(t_tskid string, t_OK, t_KO int) (b_ret bool) {
	dbopt,bret := dbbase.NewSxDB(&m_cfg.Db,"updateAbstractNum");if !bret {return}
	defer dbopt.Close()

	dbopt.Sqlcmd = "UPDATE msgAbstract SET numSentOK=numSentOK+?,numSentKO =numSentKO+? WHERE numMsg=?"
	if !dbopt.Exc(t_OK, t_KO, t_tskid) {return}
	affcted := dbopt.Affected()
	if affcted < 1 {
		util.L3E("updateAbstractNum RowsAffected()=%d no relative record(tsk=%s nOK=%d nKO=%d) update", affcted, t_tskid, t_OK, t_KO)
		return
	}

	b_ret = true
	util.L2I("updateAbstractNum OK tsk=%s affcted=%d nOK=%d nKO=%d", t_tskid, affcted, t_OK, t_KO)

	return
}
