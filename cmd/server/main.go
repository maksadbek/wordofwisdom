package main

import (
	"bufio"
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/maksadbek/wordofwisdom/hashcash"
)

var (
	addr = os.Getenv("ADDR")
)

type App struct {
	verifier *hashcash.Hashcash
	qdb      *quotedb
}

var app = App{
	verifier: hashcash.NewHashcash(time.Now, hashcash.RandStringFunc(10), 20, time.Hour*48),
	qdb:      &quotedb{db: quotes},
}

func (a *App) handleVerify(s string) (string, error) {
	if !a.verifier.Verify(s) {
		return "", errors.New("invalid token")
	}

	return a.qdb.getrandom(), nil
}

func (a *App) handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(time.Minute))
	conn.SetWriteDeadline(time.Now().Add(time.Minute))

	message, err := reader.ReadString('\n')
	if err != nil {
		log.Println("failed to read string from connection:", err)
		return
	}

	// trim right to remove the newline character
	message = message[:len(message)-1]

	log.Printf("received a challenge: %q", string(message))

	resp, err := a.handleVerify(message)
	if err != nil {
		log.Println("failed to verify token:", err)
		conn.Write([]byte(err.Error() + "\n"))
		return
	}

	log.Printf("genuine token, responding with quote: %q", resp)
	conn.Write([]byte(resp + "\n"))
}

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	log.Println("starting a tcp server at", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go app.handleRequest(conn)
	}
}

var quotes = []string{
	"The man who asks a question is a fool for a minute, the man who does not ask is a fool for life. - Confucius",
	"Do the difficult things while they are easy and do the great things while they are small. A journey of a thousand miles must begin with a single step. - Lao Tzu",
	"It is not because things are difficult that we do not dare; it is because we do not dare that things are difficult. - Seneca",
	"The only true wisdom is in knowing you know nothing. - Socrates",
	"The greatest wealth is to live content with little. - Plato",
	"In the midst of chaos, there is also opportunity. - Sun Tzu",
	"Happiness and freedom begin with one principle. Some things are within your control and some are not. - Epictetus",
}

type quotedb struct {
	db   []string
	lock sync.RWMutex
}

func (q *quotedb) getrandom() string {
	q.lock.RLock()
	defer q.lock.RUnlock()

	i := rand.Intn(len(q.db))

	return q.db[i]
}
