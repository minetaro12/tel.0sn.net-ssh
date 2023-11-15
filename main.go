package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gliderlabs/ssh"
)

var (
	hostKey *string = flag.String("f", "host.key", "SSH Host Key File")
)

func main() {
	flag.Parse()
	listen := fmt.Sprintf(":%s", getEnv("PORT", "8022"))
	count := loadCounter()

	ssh.Handle(func(s ssh.Session) {
		echoHandler(s, &count)
	})

	log.Println("SSH Server Listening on", listen, "...")
	go func() { log.Fatal(ssh.ListenAndServe(listen, nil, ssh.HostKeyFile(*hostKey))) }()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	<-ctx.Done()
	log.Println("Shutting down...")
	saveCounter(count)
}

func echoHandler(s ssh.Session, count *int) {
	log.Println("Connected:", s.RemoteAddr().String())
	*count++
	text := `
Press 'q' to exit.
--------------------------
tel.0sn.netへようこそ！
_       _   ___                         _   
| |_ ___| | / _ \ ___ _ __    _ __   ___| |_ 
| __/ _ \ || | | / __| '_ \  | '_ \ / _ \ __|
| ||  __/ || |_| \__ \ | | |_| | | |  __/ |_ 
 \__\___|_(_)___/|___/_| |_(_)_| |_|\___|\__|
--------------------------
あなたは %d 人目の訪問者です。
--------------------------
Web: https://0sn.net
--------------------------
`
	// qで終了
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			// キャンセル通知が来たら終了
			case <-ctx.Done():
				return
			default:
				buffer := make([]byte, 1)
				s.Read(buffer)
				if string(buffer) == "q" {
					s.Exit(0)
					return
				}
			}
		}
	}()

	for _, v := range fmt.Sprintf(text, *count) {
		_, err := io.WriteString(s, string(v))
		if err != nil {
			// 途中で切断された場合
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	cancel()
	log.Println("Disconnected:", s.RemoteAddr().String())
}

func loadCounter() int {
	if _, err := os.Stat("counter.txt"); os.IsNotExist(err) {
		os.WriteFile("counter.txt", []byte("0"), 0644)
	}

	data, err := os.ReadFile("counter.txt")
	if err != nil {
		log.Fatal(err)
	}

	c, err := strconv.Atoi(string(data))
	if err != nil {
		log.Fatal(err)
	}
	return c
}

func saveCounter(counter int) {
	err := os.WriteFile("counter.txt", []byte(strconv.Itoa(counter)), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v, isFound := os.LookupEnv(key); isFound {
		return v
	} else {
		return fallback
	}
}
