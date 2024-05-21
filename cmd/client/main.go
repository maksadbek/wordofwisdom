package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"net/textproto"
	"os"
	"time"

	"github.com/maksadbek/wordofwisdom/hashcash"
)

var (
	addr = os.Getenv("ADDR")
)

type App struct {
	generator *hashcash.Hashcash
}

var app = App{
	generator: hashcash.NewHashcash(time.Now, hashcash.RandStringFunc(10), 20, time.Hour*48),
}

func main() {
	if addr == "" {
		log.Fatal("please set ADDR and ID")
		return
	}

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		log.Fatal("failed to connect to server", err)
		return
	}

	conn.SetReadDeadline(time.Now().Add(time.Second * 10)) // a minute should be enough
	conn.SetWriteDeadline(time.Now().Add(time.Minute))

	defer conn.Close()

	r := textproto.NewReader(bufio.NewReader(conn))
	w := textproto.NewWriter(bufio.NewWriter(conn))

	err = w.PrintfLine("CHALLENGE")
	if err != nil {
		panic(err)
	}

	challenge, err := r.ReadLine()
	if err != nil {
		panic(err)
	}

	log.Printf("received message: %q", challenge)

	now := time.Now()
	token, err := app.generator.Generate(context.Background(), challenge)
	if err != nil {
		log.Fatal("failed to generate token", err)
	}

	log.Printf("generated a token in %v: %q", time.Since(now), token)

	err = w.PrintfLine("GETQUOTE %v", token)
	if err != nil {
		panic(err)
	}

	quote, err := r.ReadLine()
	if err != nil {
		panic(err)
	}

	log.Printf("received quote: %q", quote)
}
