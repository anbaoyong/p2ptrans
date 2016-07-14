/*
根据公司背景
客户端基于ftp下载
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
    limitrate := c.Query("limitrate")
    cutdirs := c.Query("cutdirs")
    go  wget(src,srcpath,dstpath,master,localhost,port,limitrate,cutdirs)
   c.String(200,"客户端返回")
    })
    router.Run(":12306")
}

func wget(src,srcpath,dstpath,master,localhost,port,limitrate,cutdirs string) {
    defer wg.Done()
    download := fmt.Sprintf("wget -m -r -nH -P %s --limit-rate=%s --cut-dirs=%s ftp://%s%s",dstpath,limitrate,cutdirs,src,srcpath)
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
