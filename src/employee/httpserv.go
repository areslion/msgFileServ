package employee
import(
	"net/http"
	"strconv"

	"util"
)

var m_empl sxEmp

func init(){
	http.HandleFunc("/man/getDepart",getDepartment)
	http.HandleFunc("/man/getMen",getMen)
	http.HandleFunc("/man/MenChanged",menChanged)
	http.HandleFunc("/man/search",search)
}

func getDepartment(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("%s %s",t_ask.Method,t_ask.URL)

	if t_ask.Method=="GET"{
		strPath:= t_ask.FormValue("path")
		strSep:= t_ask.FormValue("sep")
		util.L3I("%s %s",strPath,strSep)

		_,jx :=m_empl.GetLstDepat(strPath,strSep)
		t_res.Write([]byte(jx))
	}
}

func getMen(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("%s %s",t_ask.Method,t_ask.URL.Path)

	if t_ask.Method=="GET"{
		strPath:= t_ask.FormValue("path")
		strSep:= t_ask.FormValue("sep")
		strSub:= t_ask.FormValue("sub")
		util.L3I("%s %s %s",strPath,strSep,strSub)
		nsub,err := strconv.Atoi(strSub);if err!=nil{
			util.L4E("strconv.Atoi(%s) %s",strSub,err.Error())
			nsub = 1
		}
		_,bst :=m_empl.GetLstMan(strPath,strSep,nsub)
		t_res.Write(bst)
	}
}

func menChanged(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	if t_ask.Method=="POST"{
		m_empl.load()
	}
}

func search(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)
	if t_ask.Method=="GET"{
		keys:=t_ask.FormValue("keys")
		sept := t_ask.FormValue("sep")

		_,bts := m_empl.org.GetLstSearch(keys,sept)
		t_res.Write(bts)
	}
}



func StartServ(){
	m_empl.load()
	//m_empl.org.saveJson(".\\tree.json")
}