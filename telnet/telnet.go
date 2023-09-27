package telnet

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/lmj/mqtt-clients-demo/config"
)

var conns []net.Conn
var ln net.Listener
var err error

func lineCount(file *os.File) int {
	fi, _ := file.Stat()
	fileSize := fi.Size()
	chunkSize := 1024
	bytes := make([]byte, chunkSize)
	count := 0

	for i := int64(0); i < fileSize; i += int64(chunkSize) {
		n, _ := file.ReadAt(bytes, fileSize-i-int64(chunkSize))
		for j := 0; j < n; j++ {
			if bytes[j] == '\n' {
				count++
			}
		}
	}
	return count
}

func handleClient(conn net.Conn, interrupt chan os.Signal) {
	defer conn.Close()

	fileName := fmt.Sprintf(config.LOG_PATH, time.Now().Format("2006-01-02"))
	// 创建日志文件
	file, _ := os.OpenFile(fileName, os.O_RDWR, 0777)
	defer file.Close()

	// 按行读取文件并逐行发送给客户端
	scanner := bufio.NewScanner(file)
	count := lineCount(file)
	if count > 10 {
		for i := 0; i < count-10; i++ {
			scanner.Scan()
		}
	}
	for {
		select {
		case <-interrupt:
			return
		default:
			for scanner.Scan() {
				conn.Write(append(scanner.Bytes(), '\n', '\r'))
			}
			scanner = bufio.NewScanner(file)
		}
	}
}

func Start() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ln, err = net.Listen("tcp", ":23")
	if err == nil {
		fmt.Println("Telnet server started")
		for {
			conn, err := ln.Accept()
			if err == nil {
				conns = append(conns, conn)
				go handleClient(conn, interrupt)
			}
		}
	}
}

func Stop() {
	for _, conn := range conns {
		conn.Close()
	}
	ln.Close()
}
