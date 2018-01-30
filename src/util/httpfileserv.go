package util
import(
	"strings"
	"net/http"
)

func NewFileServ(t_ask *http.Request,t_res *http.ResponseWriter,t_local string){
	urlPath := t_ask.URL.Path

	L3I("NewFileServ new file dowload:"+urlPath)
	ix := strings.LastIndex(urlPath,"/");if ix==-1 {
		L3I("NewFileServ strings.LastIndex find prefix %s %d",urlPath,ix)
		return
	}
	prefix := urlPath[:ix+1]

	strx := urlPath[:ix]
	ix = strings.LastIndex(strx,"/");if ix==-1 {
		L3I("NewFileServ strings.LastIndex find folder id %s %d",urlPath,ix)
		return
	}
	folderid := strx[ix+1:]

	dirx := t_local + GetOSSeptor()+folderid
	L3I("NewFileServ start fileservr(" + prefix + " " + dirx + ")")
	staticFServ := http.StripPrefix(prefix, http.FileServer(http.Dir(dirx)))
	staticFServ.ServeHTTP(*t_res, t_ask)
}

