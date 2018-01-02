package httpserv
import (
	"log"
	"net/http"
)

var (
	staticHandler http.Handler
)


func GetlstApp(t_res http.ResponseWriter,t_ask *http.Request){
	logx("GetlstApp called")

	if t_ask.Method =="GET" {

	}
}

func logx(t_msg string){
	log.Println("fileserv  ",t_msg)
}