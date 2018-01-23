package employee
import(
	"net/http"

	"util"
)

var m_org sxOrg
var m_lstMan sxManList
func init(){
	http.HandleFunc("/man/getDepart",getDepartment)
	http.HandleFunc("/man/getMen",getMen)
}

func getDepartment(t_res http.ResponseWriter,t_ask *http.Request){
	util.L2I(t_ask.Method)

	if t_ask.Method=="GET"{
		strPath:= t_ask.FormValue("path")
		strSep:= t_ask.FormValue("sep")
		_,jx :=m_org.GetLstDepat(strPath,strSep)
		t_res.Write([]byte(jx))
	}
}

func getMen(t_res http.ResponseWriter,t_ask *http.Request){
	util.L2I(t_ask.Method)

	if t_ask.Method=="GET"{
		strPath:= t_ask.FormValue("path")
		strSep:= t_ask.FormValue("sep")
		_,bst :=m_org.GetLstMan(strPath,strSep)
		t_res.Write(bst)
	}
}

func StartServ(){
	m_lstMan.readAllMan()
	for _,itm := range m_lstMan.mapLstMan {
		m_org.insertChild(&itm)
	}
}