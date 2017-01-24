package main

import (
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mgoURL        string
	dbName        string
	cName         string
	docSleep      int
	goNum         int
	docNum        int
	docSize       int
	writeTotalNum int64
)

func init() {
	flag.StringVar(&mgoURL, "mgo_url", "mongodb://127.0.0.1:27017", "mongodb链接串")
	flag.StringVar(&dbName, "db", "test", "数据库")
	flag.StringVar(&cName, "c", "test1", "集合名")
	flag.IntVar(&docSleep, "sleep", 0, "写入文档的休眠间隔(单位毫秒，0则不休眠)")
	flag.IntVar(&goNum, "go_num", 10, "goroutine数量")
	flag.IntVar(&docNum, "doc_num", 100, "写入文档数量")
	flag.IntVar(&docSize, "doc_size", 64, "文档大小(单位B)")
}

func main() {
	flag.Parse()
	defer glog.Flush()

	session, err := mgo.DialWithTimeout(mgoURL, time.Second)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	glog.Infoln("mongodb链接建立完成")

	ticker := time.NewTicker(time.Second)
	startTime := time.Now()

	fmt.Println("开始执行...")
	go func() {
		for ct := range ticker.C {
			fmt.Printf("执行耗时：%.2fs，写入文档数：%d \r", ct.Sub(startTime).Seconds(), atomic.LoadInt64(&writeTotalNum))
		}
		fmt.Print("\n")
	}()

	docChan := make(chan []byte, docNum/10)

	go genDoc(docChan)

	wg := new(sync.WaitGroup)
	wg.Add(goNum)

	for i := 0; i < goNum; i++ {
		go writeDoc(wg, i, docChan, session)
	}

	wg.Wait()
	ticker.Stop()

	consumeTime := time.Since(startTime).Seconds()
	secTime := 1
	if consumeTime > 1 {
		secTime = int(consumeTime)
	}

	fmt.Printf("执行完成：Goroutine(%d) - 总文档数(%d) - 耗时(%.2fs) - 写入成功数量(%d) - 每秒写入数量(%d) - 失败数量(%d) \n",
		goNum, docNum, consumeTime, writeTotalNum, writeTotalNum/int64(secTime), docNum-int(writeTotalNum))
}

func genDoc(writeChan chan<- []byte) {
	docs := make([]byte, docSize)
	for i := 0; i < docSize; i++ {
		docs[i] = 'a'
	}

	n := 0
	for n < docNum {
		writeChan <- docs
		n++
	}

	close(writeChan)
	docs = nil

	glog.Infoln("文档生成完成，等待写入")
}

func writeDoc(wg *sync.WaitGroup, num int, docChan <-chan []byte, session *mgo.Session) {
	defer wg.Done()
	i := 0
	for doc := range docChan {
		sess := session.Clone()
		err := sess.DB(dbName).C(cName).Insert(bson.M{"text": doc})
		sess.Close()
		if err != nil {
			glog.Errorf("G[%d]插入文档发生异常:%s \n", num, err.Error())
			continue
		}
		atomic.AddInt64(&writeTotalNum, 1)
		i++
		if docSleep > 0 {
			time.Sleep(time.Millisecond * time.Duration(docSleep))
		}
	}
	glog.Infof("G[%d]已经关闭,写入文档数[%d] \n", num, i)
}
