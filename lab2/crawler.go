package main
 
import (
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "errors"
    "strings"
    "time"

    "os"
	"syscall"
)

type Crawler struct{
    routChan chan byte
    urlChan chan string
    filterChan chan string
    filter map[string]bool
    field string
    closeChan chan byte
    start time.Time
}

func NewCrawler() *Crawler{
    return &Crawler{
        routChan:make(chan byte,40),
        urlChan:make(chan string,10000),
        filterChan:make(chan string,100000),
        filter:make(map[string]bool),
        closeChan:make(chan byte),
        start:time.Now(),
    }
}

func(c * Crawler)demon(){
    dida := 0
    for{
        c.showState()
        if len(c.routChan) == 0{
            dida ++
        }else{
            dida = 0
        }
        if dida >=5{
            close(c.closeChan)
            return
        }
        time.Sleep(5 * time.Second)
    }
}

func(c * Crawler)filterUrl(){
    filename := "./"+c.field + ".txt"
    f,err := os.Create(filename)
    if err != nil{
        fmt.Println(err)
        close(c.closeChan)
        return
    }
    defer f.Close()
    for{
        select{
            case url := <- c.filterChan:
                if _,ok := c.filter[url]; ok{
                    continue
                }
                c.filter[url] = true
                f.Write([]byte(url+"\n"))
                f.Close()
                f,_ = os.OpenFile(filename,syscall.O_APPEND,777)
                c.urlChan <- url
            
            case <-c.closeChan:
                return
        }
    }
}

func(c* Crawler)showState(){
    fmt.Println("/***************************************/")
    fmt.Println("Crawler in field:",c.field)
    fmt.Println("Routine Now:",len(c.routChan))
    fmt.Println("Url now found:",len(c.filter))
    fmt.Println("Url num in filter channel:",len(c.filterChan))
    fmt.Println("Url num in waiting channel:",len(c.urlChan))
    fmt.Println("Crawler has run:",time.Now().Sub(c.start))
    fmt.Println("/***************************************/")
}

func (c * Crawler)Get(url string) (content string, err error){
    timeout := time.Duration(3 * time.Second)
    client := http.Client{
        Timeout:timeout,
    }
    resp,err := client.Get(url)
    
    if err != nil{
        return 
    }
    defer resp.Body.Close()
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil{
        return
    }
    if resp.StatusCode != 200{
        str := fmt.Sprintf("StatusCode is %d.\n",resp.StatusCode)
        err = errors.New(str)
        return
    }
    content = string(data)
    return
}

func (c *Crawler)Analyser(ctx string){
    reURL := regexp.MustCompile("<a.*?href=\"(.*?)\"")
    matches := reURL.FindAllStringSubmatch(ctx,10000)
    for _,url := range matches{
        if strings.Contains(url[1],c.field){
                c.filterChan <- url[1]
        }
    }
}

func (c *Crawler)Run(seed,field string){
    c.field = field
    ctx, err := c.Get(seed)
    if err != nil{
        fmt.Println(err)
        return
    }
    c.filter[seed] = true
    go c.filterUrl()
    c.Analyser(ctx)
    go c.demon()
    count := 0
    var url string
    for {
        select{
            case url = <- c.urlChan:
                c.routChan <- 1
                count ++
                go func(id int){
                     
                    html,err1 := c.Get(url)
                    if err1 != nil{
                        fmt.Println(err1)
                        <-c.routChan
                        return 
                    }
                    c.Analyser(html)
                    <-c.routChan
                }(count)
                case <-c.closeChan:
                    fmt.Println("Crawler Done!")
                    c.showState()
                    return 

        }
    }
}

func main(){
    c := NewCrawler()
    c.Run("http://www.bupt.edu.cn/","bupt.edu.cn")
}