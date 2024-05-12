package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/maksadbek/wordofwisdom/proofofwork"
	"github.com/maksadbek/wordofwisdom/proofofwork/hashcash"
)

var (
	addr = os.Getenv("ADDR")
	id   = os.Getenv("ID")
)

type App struct {
	generator proofofwork.Generator
}

var app = App{
	generator: hashcash.NewHashcash(time.Now, hashcash.RandStringFunc(10), 20, time.Hour*48),
}

func main() {
	if addr == "" || id == "" {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	now := time.Now()
	token, err := app.generator.Generate(ctx, id)
	if err != nil {
		log.Fatal("failed to generate token", err)
	}

	log.Printf("generated a token in %v: %q", time.Since(now), token)

	_, err = io.WriteString(conn, token+"\n")
	if err != nil {
		log.Fatal("failed to send token", err)
	}

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("failed to read from conn:", err)
		return
	}

	message = message[:len(message)-1]
	log.Printf("received message: %q", message)
}
