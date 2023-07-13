package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func main() {
	tFlag := flag.Uint64("timeout", 10, "connection timeout")
	flag.Parse()
	t := *tFlag
	if flag.NArg() != 2 {
		log.Fatal("host and port not provided (usage: go-telnet [--timeout=5] <host> <port>)")
	}
	client := New(flag.Args()[0], flag.Args()[1], t)
	if err := client.Run(); err != nil {
		log.Fatal(err)
	}
}

// Client contains address to connect to and timeout for connection establishment
type Client struct {
	Host    string
	Port    string
	Timeout time.Duration
}

// New creates an instance of Client provided with host, port and timeout in seconds and returns a pointer to it
func New(host, port string, t uint64) *Client {
	return &Client{
		Host:    host,
		Port:    port,
		Timeout: time.Duration(t) * time.Second,
	}
}

// Run connects to telnet server, so client can send and recieve messages
func (c *Client) Run() error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", c.Host, c.Port), c.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go send(ctx, stop, conn)
	go recieve(ctx, stop, conn)
	<-ctx.Done()
	return nil
}

func send(ctx context.Context, stop context.CancelFunc, conn net.Conn) {
	r := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := r.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					stop()
					conn.Close()
					return
				}
				log.Println(err)
				continue
			}
			if _, err = conn.Write(msg); err != nil {
				log.Println(err)
			}
		}
	}
}

func recieve(ctx context.Context, stop context.CancelFunc, conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if msg, err := r.ReadString('\n'); err != nil {
				if err == io.EOF {
					stop()
					conn.Close()
					return
				}
				log.Println(err)
			} else {
				fmt.Println(msg)
			}
		}
	}
}
