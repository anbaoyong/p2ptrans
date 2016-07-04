package main

import (
  "github.com/gin-gonic/gin"
  "fmt"
  "strconv"
//  "github.com/BurntSushi/toml"
  "github.com/parnurzeal/gorequest"
//  "io/ioutil"
  "io"
  "strings"
  "bufio"
  "os"
  "time"
  "flag"
)
//定义slice
var (
SuccList []string //成功列表
FailList []string //失败列表
DstList []string //目的列表 
//p2p Config
ListenPort string
Port *string
Master *string
Srcpath *string
Dstpath *string
Filename *string
Timeout *int
)
var succ_ch = make(chan string,10000)
/*
type Config struct {
    Masterid string
    Timeout int //暂未使用
    Port string
}
*/
//读取文件，返回行 列表
func Cat(filename string) []string {
   f,err := os.Open(filename) 
   defer f.Close()
   if err != nil {
        panic(err)
   } 
   reader := bufio.NewReader(f)
   for {
        line,err := reader.ReadString('\n')
        if err == io.EOF { break }
        line = strings.TrimSpace(line)
        if line != "" {
            DstList = append(DstList,line)
        }
   }
   return DstList
}
//读取配置文件
/*
func Conf() *Config {
    bs, _ := ioutil.ReadFile("p2p.conf")
    if _, err := toml.Decode(string(bs), &p2p); err != nil {
        fmt.Println("decode config file failed: %s\n", err.Error())
    }
    return &p2p
}
*/
//接收客户端返回的信息
func accept() {
  router := gin.Default()
  router.GET("/p2p",func(c *gin.Context){
  status,_ := strconv.ParseBool(c.Query("status")) //判断客户端是否下载成功
  host := c.Query("host")
  src := c.Query("src") //接受客户端返回的数据源
  //fmt.Println(status,host)
  if status {
      succ_ch <- host
      succ_ch <- src
      SuccList = append(SuccList,host)
     // fmt.Println(SuccList,FailList)
  } else {
      fmt.Println("客户端下载失败：",host)
      succ_ch <- src
  }
  })
  router.Run(ListenPort)
  //router.Run(":%s",Conf().Port)
}
//request向客户端发送下载任务
func request(master,port,src,dst,srcpath,dstpath string) {
    req := gorequest.New()
    url := fmt.Sprintf("http://%s:12306/client?port=%s&src=%s&srcpath=%s&dstpath=%s&master=%s&localhost=%s",dst,port,src,srcpath,dstpath,master,dst)
    fmt.Println(url)
    //_, _, errs := req.Get(url).Timeout(10*time.Second).End() //设置url超时时间为10妙
    _, _, errs := req.Get(url).End() //设置url超时时间为10妙
    if len(errs) != 0 {
        fmt.Println("请求客户端失败:",dst)
        succ_ch <- src
    }
}
//主程序执行分发控制
func handle(dstlist []string) {
    index := 0
    lendst := len(dstlist)
    for {
        select {
            case succhost := <- succ_ch:
                if index < lendst {
                    request(*Master,*Port,succhost,dstlist[index],*Srcpath,*Dstpath)
                    index = index + 1
                } else {
                    time.Sleep(time.Duration(*Timeout)*time.Second)
                    if lendst != len(SuccList) {
                        FailList = Set(dstlist,SuccList)
                        fmt.Printf("成功机器数为%d，机器列表为，机器列表为%v\n",len(SuccList),SuccList)
                        fmt.Printf("失败机器数为%d，机器列表为，机器列表为%v\n",len(FailList),FailList)
                    } else {
                        fmt.Printf("全部成功，成功机器数为%d，机器列表为，机器列表为%v\n",len(SuccList),SuccList)
                    }
                    fmt.Println("==============执行完毕=================")
                    os.Exit(0)
                } 
        }
    } 
}
//初始化
func init() {
//解析command line参数
    flag.Usage = func() {
        fmt.Println("Usage: <-m host> [-f file] [-s srcpath] [-d dstpath] [-p port] [-t timeout]")
        flag.PrintDefaults()
    }
    Filename = flag.String("f","ip.txt","File that contains the target machine")
    Srcpath = flag.String("s","/home/xiaoju","Data source path")
    Dstpath = flag.String("d","/home/xiaoju","Data destination path")
    Port =  flag.String("p","12306","Listen port")
    Master = flag.String("m","","ip or host name of the master")
    Timeout = flag.Int("t",1800,"If the master does not receive the return value within the specified time , that the transmission fails")
    flag.Parse()
    ListenPort = fmt.Sprintf(":%s",*Port)
    succ_ch <- *Master //将master加入channel
}
//两个slice取差集
func Set(one, two []string) []string {
    x := []string{}
    if len(two) != 0 {
    for _, v := range one {
        for k, vv := range two {
            if v == vv {
                break
            }
            if k == len(two)-1 {
                x = append(x, v)
            }
        }
    }
    } else {
        return one
    }
    return x
}
func main() {
    if *Master == "" {
        flag.Usage()
        os.Exit(2)
    }
    go accept()
    handle(Cat(*Filename))
}
