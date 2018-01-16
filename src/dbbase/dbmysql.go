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
