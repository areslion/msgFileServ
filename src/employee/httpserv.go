package employee
import(
	"net/http"
	"strconv"

	"util"
	"io/ioutil"
)

var m_empl sxEmp
var m_group sxGroup

func init(){
	http.HandleFunc("/man/getDepart",getDepartment)
	http.HandleFunc("/man/getMen",getMen)
	http.HandleFunc("/man/group",manGroup)
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

func manGroup(t_res http.ResponseWriter,t_ask *http.Request){
	util.L3I("%s %s",t_ask.Method,t_ask.URL.Path)

	if t_ask.Method=="GET"{
		jsx,bts:=m_group.GetGroup(cst_sepstd)
		jsx = jsx
		//t_res.Write([]byte(jsx))
		util.L3I("jsx=%d,bts=%d",len(jsx),len(bts));
		t_res.Write(bts)
	} else if t_ask.Method=="POST"{
		bts,err := ioutil.ReadAll(t_ask.Body);if err!=nil{
			util.L4E("manGroup "+err.Error())
		}
		m_group.saveGroup(bts)
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