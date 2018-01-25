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
	http.HandleFunc("/man/MenChanged",menChanged)
}

func getDepartment(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	if t_ask.Method=="GET"{
		strPath:= t_ask.FormValue("path")
		strSep:= t_ask.FormValue("sep")

		_,jx :=m_org.GetLstDepat(strPath,strSep)
		t_res.Write([]byte(jx))
	}
}

func getMen(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	if t_ask.Method=="GET"{
		strPath:= t_ask.FormValue("path")
		strSep:= t_ask.FormValue("sep")
		_,bst :=m_org.GetLstMan(strPath,strSep)
		t_res.Write(bst)
	}
}

func menChanged(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	if t_ask.Method=="POST"{
		loadOrg()
	}
}

func loadOrg(){
	m_lstMan.lstMan = m_lstMan.lstMan[0:0]
	m_org.clear()
	m_lstMan.readAllMan()
	for _,itm := range m_lstMan.lstMan {
		m_org.insertChild(&itm)
	}

	util.L3I("man info has bee re-load,num=%d",len(m_lstMan.lstMan))
}

func StartServ(){
	loadOrg()
}