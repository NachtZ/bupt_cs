package main
import (
    "fmt"
    "os"
    "io/ioutil"
    "strings"
    "time"
)
type PatTrie struct{
    childs []*PatTrie
    val *string
    num int
}

type Pattrie interface{
    insert_node(key string, bin []int,tmp []int)
    find(key string, bin []int, tmp []int) bool
    init(key string ,bin []int) 
    
}

func NewPatTrie(num int) *PatTrie{
    p := &PatTrie{
        childs:make([]*PatTrie,2),
        val:nil,
        num:num,
    }
    return p
}

var bincode [][]int

func compare(bin1,bin2 []int)int{
    i := 0
    len1 := 0
    if len(bin1)>len(bin2){
        len1 = len(bin2)
    }else{
        len1 = len(bin1)
    }
    for i=0;i<len1 && bin1[i] == bin2[i];i++{
    }
    if i == len1 && bin1[i]!=bin2[i]{
        return -1
    }
    return i
}
func(p *PatTrie) init(key string, bin []int){
    strtobin(key,bin)
    p.num = 0
    p.childs[bin[0]] = NewPatTrie(-1)
    p.childs[bin[0]].val = &key
}

func(p *PatTrie)insert_node(key string,bin []int,tmp []int){
    //bin := make([]int,0,len(key)*6)
    len1 := strtobin(key,bin)
    node,pre := p,p
    for{
        i := node.num
        if i > len1 || i == -1 || node.childs[bin[i]] == nil{
            if i <= len1 && i != -1{
                node.childs[bin[i]] = NewPatTrie(-1)
                node.childs[bin[i]].val = &key
                return
            }
            strtobin(*node.val,tmp)
            t := compare(tmp,bin)
            if t == -1{
                return
            }
            if pre.num>t{
                node = p
                for t>node.num{
                    pre = node
                    node = node.childs[bin[node.num]]
                }
            }
            pre.childs[bin[pre.num]] = NewPatTrie(t)
            pre.childs[bin[pre.num]].val = &key
            pre.childs[bin[pre.num]].childs[bin[t]] = NewPatTrie(-1)
            pre.childs[bin[pre.num]].childs[bin[t]].val = &key
            pre.childs[bin[pre.num]].childs[1-bin[t]] = node
            return
        }else{
            pre = node
            node = node.childs[bin[i]]
        }
    }
}

func(p *PatTrie)find(key string,bin []int,tmp []int) bool{
    strtobin(key,bin)
    node:=p
    for node != nil{
        if node.num == -1{
            if node.val == nil{
                return false
            }
            fmt.Println(*node.val)
            return *node.val == key
        }
        node = node.childs[bin[node.num]]
    }
    return false
}

func printf(){
    for i:=0;i<128;i++{
        fmt.Printf("%d:",i)
        fmt.Println(bincode[i])
    }
}

func strtobin(key string, bin []int) int{
    lenth := 0
    j :=0
    if len(key) == 0{
        return 0
    }
    strlen := len(key)
    for i := strlen -1;i>=0;i--{
        for k:=0;k<6;k++{
            bin[j] = bincode[key[i]][k]
            j++
        }
        lenth += 6
    }
    for k:=0;k<6;k++{
        bin[j] = bincode[40][k]
        j++
    }
    lenth +=6
    return lenth
}

func initFunc(){
    k := -1
    bincode = make([][]int,256)
    for i:=0;i<123;i++{
        tmp := make([]int,6)
        if i<='z' && i >= 'a'{
            k = i -'a'
        }else{
            if i<='Z' && i>='A'{
                k = i - 'A'
            }else{
                if i <= '9' && i>= '0'{
                    k = i -'0' +26
                }
            }
        }
        if k == -1{
            switch i{
                case '.':
                    k = 36
                case '@':
                    k = 37
                case '-':
                    k = 38
                case '_':
                    k = 39
                case 10:
                    k = 40
                case 13:
                    k = 40
                default:
                    bincode[i] = []int{1,1,1,1,1,1}
                    continue
            }
        }
        for j:=0;j<6;j++{
            if( uint16(k)>>uint16(j) & 1 )>= 1{
                tmp[j] = 1
            }else{
                tmp[j] = 0
            }
        }
        bincode[i] = tmp
    }


}

func run(mail,check string){
    fi,err := os.Open(mail)
    if err!=nil{
        panic(err)
    }
    defer fi.Close()
    fd,err := ioutil.ReadAll(fi)
    str := string(fd)
    fc,err := os.Open(check)
    if err !=nil{
        panic(err)
    }
    defer fc.Close()
    fd1,err := ioutil.ReadAll(fc)
    checks := string(fd1)
    trie := NewPatTrie(-1)
    count :=0
    bin := make([]int,321*6)
    tmp := make([]int,321*6)
    initFunc()
    //printf()
    trie.init("tsetInit",bin)
    for _,s := range strings.Split(str,"\n"){
        if s == ""{
            continue
        }
        trie.insert_node(s,bin,tmp)
    }
    fmt.Println("0.1z287khq3y@listserv.workplacetoolbox.nerpveadlvvbqfget",":",trie.find("0.1z287khq3y@listserv.workplacetoolbox.nerpveadlvvbqfget",bin,tmp))
    for _,s:= range strings.Split(checks,"\r\n"){
        if s == ""{
            continue
        }
        t := trie.find(s,bin,tmp)
        fmt.Println(s,":",t)
        if t {
            count ++
        }
    }
    fmt.Println(count)
}

func main(){
    start := time.Now()
    run("F:/emaillist.dat","F:/checklist.dat")
    end := time.Now()
    fmt.Println(end.Sub(start))
}