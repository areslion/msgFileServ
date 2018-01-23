package dbbase

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)
import (
	"fmt"
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
	tag string
}

func NewSxDB(t_cfg *util.SxCfg_db, t_tag string) (r_new *sxDB, b_ret bool) {
	r_new = &sxDB{cfg: t_cfg, tag: t_tag}
	b_ret = r_new.open()
	return
}
func (p *sxDB) Affected()(r_afc int64) {
	var err error
	r_afc, err = p.res.RowsAffected()
	if err != nil {
		p.logF("updateSendStatus res.RowsAffected() " + err.Error())
	}
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
	if !p.PrePare() {
		return
	}

	p.res, err = p.smt.Exec(args...)
	if err != nil {
		p.logF("Exc Fail to Query(args) " + err.Error())
		return
	}

	b_ret = true
	return
}
func (p *sxDB) ExcAlone(args ...interface{}) (b_ret bool) {
	var err error
	p.res, err = p.smt.Exec(args...)
	if err != nil {
		p.logF("Exc Fail to Query(args) " + err.Error())
		return
	}

	b_ret = true
	return
}

func (p *sxDB) open() (b_ret bool) {
	var err error
	p.cnn, err = sql.Open("mysql", p.cfg.GetCntStr())
	if err != nil {
		p.logF("open() Fail to open db " + err.Error() + " " + p.cfg.GetCntStr())
		return
	}
	b_ret = true
	return
}
func (p *sxDB) Query(args ...interface{}) (b_ret bool) {
	var err error
	if !p.PrePare() {
		return
	}

	p.Rows, err = p.smt.Query(args...)
	if err != nil {
		p.logF("Query() Fail to Query(args) " + err.Error())
		return
	}

	b_ret = true
	return
}
func (p *sxDB) PrePare() (b_ret bool) {
	p.closeRes()

	var err error
	p.smt, err = p.cnn.Prepare(p.Sqlcmd)
	if err != nil {
		p.logF("prePare() Fail to Prepare(%s) %s", p.Sqlcmd, err.Error())
		return
	}
	b_ret = true
	return
}
func (p *sxDB) logF(t_fmt string,v...interface{}){
	fmtx := fmt.Sprintf("sxDB(%s) %s",p.tag,t_fmt)
	util.L5Fx(5,2,fmtx,v...)
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

func Open(cfg *util.SxCfgAll) (r_cnt *sql.DB, r_res bool) {
	m_cfgdb = cfg

	cnn, err := sql.Open("mysql", m_cfgdb.Db.GetCntStr())
	if err != nil {
		util.L5F("Fail to open db " + err.Error() + " " + m_cfgdb.Db.GetCntStr())
		return
	}
	r_cnt = cnn
	r_res = true
	return
}
func Close(t_cnn *sql.DB) {
	t_cnn.Close()
}


