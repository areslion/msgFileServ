package util
import(
	"log"
	"os"
	"fmt"
	"github.com/satori/go.uuid"
)
const (
	Cst_sept = "/"
)


func SaveFileBytes(filename string,buf [] byte) (r_path string, b_ret bool){
	t,err :=os.Create(filename)
	if(err!=nil){
		log.Fatal("saveFile 创建服务端文件失败",err.Error())
		return "", false
	}
	defer t.Close()

	log.Println("成功创本地文件 ",filename)		
	if _,err := t.Write(buf);err!=nil{
		log.Fatal("saveFile 存储文件失败 ",err.Error())		
		return "" ,false
	}

	log.Println("成功保存文件:",filename)
	return filename,true
}

func Guid()(r_guid string){
	u1 := uuid.NewV4()
	
	// log.Println("-------------x-",u1,"-x----------")
	// fmt.Printf("-----y-%s-y-------",u1)	

	//strRet := string(u1[0:len(u1)])
	strRet := fmt.Sprintf("%s",u1)
	return strRet
}