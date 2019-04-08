package db

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"gocn/message"
	"log"
	"os"
	"time"
)

var markdownFd *os.File
var bufferChain chan string
var markdownBufferChain chan string

var (
	Db *badger.DB
)

func Run() {
	for {
		select {
		case msg := <-bufferChain:
			write(msg)
		case m := <-markdownBufferChain:
			writeMarkdown(m)
		default:
			time.Sleep(time.Second * 1)
		}
	}
}

func init() {

	bufferChain = make(chan string, 1000)
	markdownBufferChain = make(chan string, 1000)

	opts := badger.DefaultOptions
	opts.Dir = "./db/data"
	opts.ValueDir = "./db/data"
	var err error
	Db, err = badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	markdownFd, err = os.OpenFile("./daily/golang-daily.md", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err.Error())
	}
}

func Close() {
	_ = Db.Close()
}

//检测是否已经存储
func HasSave(title string) bool {
	saved := false
	err := Db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(title))
		if err == nil {
			saved = true
			fmt.Println(title)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return saved
}

func writeMarkdown(msg string) {
	_, _ = fmt.Fprintln(markdownFd, "\n"+msg)
	_ = markdownFd.Sync()
}

func write(msgStr string) {
	var msg message.Message
	_ = json.Unmarshal([]byte(msgStr), &msg)
	_ = Db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(msg.DailyTitle), []byte(message.BuildMessage(msg)))
		if err != nil {
			log.Fatal(err)
		}
		return err
	})
}

func Push(msg string) {
	bufferChain <- msg
}

func PushMarkdown(msg string) {
	markdownBufferChain <- msg
}
