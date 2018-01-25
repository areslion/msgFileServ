package employee
import(
	"net/http"

	"util"
)

var m_empl sxEmp

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

		_,jx :=m_empl.GetLstDepat(strPath,strSep)
		t_res.Write([]byte(jx))
	}
}

func getMen(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	if t_ask.Method=="GET"{
		strPath:= t_ask.FormValue("path")
		strSep:= t_ask.FormValue("sep")
		_,bst :=m_empl.GetLstMan(strPath,strSep)
		t_res.Write(bst)
	}
}

func menChanged(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I(t_ask.Method)

	if t_ask.Method=="POST"{
		m_empl.load()
	}
}



func StartServ(){
	m_empl.load()
}