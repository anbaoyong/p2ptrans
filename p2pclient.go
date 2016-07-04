/*
根据公司背景
客户端基于ftp下载
下载限速为80M，且不能修改
cut-dirs默认是2，且不能修改
此服务下载完成后会自动退出
*/

package main
import (
  "github.com/gin-gonic/gin"
  "github.com/parnurzeal/gorequest"
  "fmt"
  "os/exec"
  "sync"
)
var ( 
Url string
Status string
wg sync.WaitGroup
)

func main() {
    wg.Add(1)
    go accept()
    wg.Wait()
    request := gorequest.New()
    _, _,errs := request.Get(Url).End()
    if len(errs) != 0 {
        fmt.Println("访问master失败")
    }
}
//接收master请求
func accept() {
  router := gin.Default()
  router.GET("/client",func(c *gin.Context){
    src := c.Query("src")
    master := c.Query("master")
    srcpath := c.Query("srcpath")
    dstpath := c.Query("dstpath")
    localhost := c.Query("localhost")
    port := c.Query("port")
    go  wget(src,srcpath,dstpath,master,localhost,port)
   c.String(200,"客户端返回")
    })
    router.Run(":12306")
}

func wget(src,srcpath,dstpath,master,localhost,port string) {
    defer wg.Done()
    download := fmt.Sprintf("wget -m -r -nH -P %s --limit-rate=80m --cut-dirs=2 ftp://%s%s",dstpath,src,srcpath)
    fmt.Println(download)
   _,err := exec.Command("sh","-c",download).CombinedOutput()
    if err != nil {
        Status = "false"
    } else {
        Status = "true"
    }
    fmt.Println(src,master,dstpath,localhost)
    Url = fmt.Sprintf("http://%s:%s/p2p?host=%s&status=%s&src=%s",master,port,localhost,Status,src)

}
