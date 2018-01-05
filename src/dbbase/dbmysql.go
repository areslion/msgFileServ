package dbbase

import(	
	"database/sql"	
	_ "github.com/go-sql-driver/mysql"
	"log"
)
import(
	"util"
)


type SCfg struct{
	hostIP string
	usrname string
	password string
	dbname string
	charset string
}
func (p* SCfg)Init(t_ip,t_usr,t_pwd,t_db,t_cset string){
	p.hostIP = t_ip
	p.usrname = t_usr
	p.password = t_pwd
	p.dbname = t_db
	p.charset = t_cset
}
func (p* SCfg)GetCntStr()string{
	cntstr := p.usrname+":"+p.password+"@tcp("+p.hostIP+")/"
	cntstr += p.dbname+"?charset="+p.charset

	return cntstr
}


var (
	m_cfgdb *util.SxCfgAll
	m_dbcnt	*sql.DB
)

func Tstmysql() bool {

	var err error
	strcnt := m_cfgdb.Db.GetCntStr()
	m_dbcnt,err = sql.Open("mysql",strcnt)
	if(err!=nil){
		logx("Fail to open db "+err.Error())
		return false
	}
	defer m_dbcnt.Close()

	rows,err := m_dbcnt.Query("SELECT numDev from terDevBasicInfo")
	if(err!=nil){
		logx("Fail to select data "+err.Error())
	}
	showRows(rows)

	return true
}

func Open(cfg *util.SxCfgAll)(r_cnt *sql.DB,r_res bool){
	m_cfgdb = cfg

	var err error
	m_dbcnt,err = sql.Open("mysql",m_cfgdb.Db.GetCntStr())
	if(err!=nil){
		logx("Fail to open db "+err.Error()+" "+m_cfgdb.Db.GetCntStr())
		return nil,false
	}
	return m_dbcnt,true;
}
func Close(){
	m_dbcnt.Close()
}


func logx(t_msg string){
	log.Println("dbmysql  "+t_msg)
}

func showRows(t_rows *sql.Rows){
	var numDev string

	for t_rows.Next(){
		
		err := t_rows.Scan(&numDev)
		if(err==nil) {
			logx("Res:"+numDev)
		}
	}
}
