package dbbase

import(	
	"database/sql"	
	_ "github.com/go-sql-driver/mysql"
	"log"
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
	m_cfgdb SCfg
	m_dbcnt	*sql.DB
)

func Tstmysql() bool {

	var err error
	strcnt := initCfg()
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

func Open(cfg *SCfg)(r_cnt *sql.DB,r_res bool){
	m_cfgdb = *cfg

	var err error
	m_dbcnt,err = sql.Open("mysql",m_cfgdb.GetCntStr())
	if(err!=nil){
		logx("Fail to open db "+err.Error()+" "+m_cfgdb.GetCntStr())
		return nil,false
	}
	return m_dbcnt,true;
}
func Close(){
	m_dbcnt.Close()
}

func initCfg()string {	
	m_cfgdb.charset = "utf8"
	m_cfgdb.dbname = "deskSafe"
	m_cfgdb.hostIP = "10.20.10.101"
	m_cfgdb.usrname = "root"
	m_cfgdb.password = "123456"

	cntstr := m_cfgdb.usrname+":"+m_cfgdb.password+"@tcp("+m_cfgdb.hostIP+")/"
	cntstr += m_cfgdb.dbname+"?charset="+m_cfgdb.charset

	return cntstr
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
