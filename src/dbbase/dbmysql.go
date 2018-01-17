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

type sxDB struct {
	cnn    *sql.DB
	Rows   *sql.Rows
	smt    *sql.Stmt
	res    sql.Result
	cfg    *util.SxCfg_db
	Sqlcmd string
}

func NewSxDB(t_cfg *util.SxCfg_db, t_sql string) (r_new *sxDB, b_ret bool) {
	r_new = &sxDB{cfg: t_cfg, Sqlcmd: t_sql}
	b_ret = r_new.open()
	return
}
func (p *sxDB) Close() {
	p.closeRes()
	p.cnn.Close()
}
func (p *sxDB) closeRes() {
	if p!=nil&&p.Rows!=nil {p.Rows.Close()}
	if p!=nil&&p.smt!=nil {p.smt.Close()}
}
func (p *sxDB) Exc(args ...interface{}) (b_ret bool) {
	var err error
	if !p.prePare() {
		return
	}

	p.res, err = p.smt.Exec(args)
	if err != nil {
		util.L4F("SxDB Fail to Query(args) " + err.Error())
		return
	}

	b_ret = true
	return
}
func (p *sxDB) open() (b_ret bool) {
	var err error
	p.cnn, err = sql.Open("mysql", p.cfg.GetCntStr())
	if err != nil {
		util.L4F("SxDB Fail to open db " + err.Error() + " " + p.cfg.GetCntStr())
		return
	}
	b_ret = true
	return
}
func (p *sxDB) Query(args ...interface{}) (b_ret bool) {
	var err error
	if !p.prePare() {
		return
	}

	p.Rows, err = p.smt.Query(args...)
	if err != nil {
		util.L4F("SxDB Fail to Query(args) " + err.Error())
		return
	}

	b_ret = true
	return
}
func (p *sxDB) prePare() (b_ret bool) {
	p.closeRes()

	var err error
	p.smt, err = p.cnn.Prepare(p.Sqlcmd)
	if err != nil {
		util.L4F("SxDB Fail to Prepare(%s) %s", p.Sqlcmd, err.Error())
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
