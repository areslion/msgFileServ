package dbbase

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)
import (
	"util"
)

type SCfg struct {
	hostIP   string
	usrname  string
	password string
	dbname   string
	charset  string
}

type SxDB struct{
	Cnn *sql.DB
	Rows *sql.Rows
	Smt *sql.Stmt
	Res sql.Result
	Cfg *util.SxCfg_db
	Sqlcmd string
}
func (p *SxDB)Close(){
	p.CloseRes()
	p.Cnn.Close()
}
func (p *SxDB)CloseRes(){
	p.Rows.Close()
	p.Smt.Close()
}
func NewSxDB(t_cfg *util.SxCfg_db)(r_new *SxDB){
	return &SxDB{Cfg:t_cfg}
}
func (p *SxDB)Open()(b_ret bool){
	var err error
	p.Cnn, err = sql.Open("mysql", p.Cfg.GetCntStr())
	if err != nil {
		util.L4F("SxDB Fail to open db " + err.Error() + " " + p.Cfg.GetCntStr())
		return
	}
	b_ret = true
	return
}
func (p *SxDB)Query(args ...interface{})(b_ret bool){
	var err error
	p.Rows,err = p.Smt.Query(args);if err!=nil {
		util.L4F("SxDB Fail to Query(args) "+err.Error())
		return
	}

	b_ret = true
	return
}
func (p *SxDB)Exc(args ...interface{})(b_ret bool){
	var err error
	p.Res,err = p.Smt.Exec(args);if err!=nil {
		util.L4F("SxDB Fail to Query(args) "+err.Error())
		return
	}

	b_ret = true
	return
}
func (p *SxDB)PrePare()(b_ret bool){
	p.CloseRes()

	var err error
	p.Smt,err = p.Cnn.Prepare(p.Sqlcmd);if err!=nil{
		util.L4F("SxDB Fail to Prepare(%s) %s",p.Sqlcmd,err.Error())
		return
	}
	b_ret = true
	return
}

func (p *SCfg) Init(t_ip, t_usr, t_pwd, t_db, t_cset string) {
	p.hostIP = t_ip
	p.usrname = t_usr
	p.password = t_pwd
	p.dbname = t_db
	p.charset = t_cset
}
func (p *SCfg) GetCntStr() string {
	cntstr := p.usrname + ":" + p.password + "@tcp(" + p.hostIP + ")/"
	cntstr += p.dbname + "?charset=" + p.charset

	return cntstr
}

var (
	m_cfgdb *util.SxCfgAll
)

func Tstmysql() bool {

	var err error
	strcnt := m_cfgdb.Db.GetCntStr()
	cnn, err := sql.Open("mysql", strcnt)
	if err != nil {
		util.L3E("Fail to open db " + err.Error())
		return false
	}
	defer cnn.Close()

	rows, err := cnn.Query("SELECT numDev from terDevBasicInfo")
	if err != nil {
		util.L3E("Fail to select data " + err.Error())
	}
	showRows(rows)

	return true
}

func Open(cfg *util.SxCfgAll) (r_cnt *sql.DB, r_res bool) {
	m_cfgdb = cfg

	cnn, err := sql.Open("mysql", m_cfgdb.Db.GetCntStr())
	if err != nil {
		util.L4F("Fail to open db " + err.Error() + " " + m_cfgdb.Db.GetCntStr())
		return
	}
	r_cnt = cnn
	r_res = true
	return
}
func Close(t_cnn *sql.DB) {
	t_cnn.Close()
}

func showRows(t_rows *sql.Rows) {
	var numDev string

	for t_rows.Next() {

		err := t_rows.Scan(&numDev)
		if err == nil {
			util.L1T("Res:" + numDev)
		}
	}
}
