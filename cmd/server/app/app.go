package app

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"net/textproto"
	"strings"
	"time"

	"github.com/maksadbek/wordofwisdom/hashcash"
	"github.com/maksadbek/wordofwisdom/internal/command"
	"github.com/maksadbek/wordofwisdom/internal/kvstore"
)

type App struct {
	commands *command.Factory

	listener net.Listener

	shutdown chan struct{}
}

func NewApp() *App {
	hc := hashcash.NewHashcash(time.Now, hashcash.RandStringFunc(10), 20, time.Hour*48)

	quotesdb := kvstore.NewStore[string, string]()
	quotesdb.Set("0", "The man who asks a question is a fool for a minute, the man who does not ask is a fool for life. - Confucius")
	quotesdb.Set("1", "Do the difficult things while they are easy and do the great things while they are small. A journey of a thousand miles must begin with a single step. - Lao Tzu")
	quotesdb.Set("2", "It is not because things are difficult that we do not dare; it is because we do not dare that things are difficult. - Seneca")
	quotesdb.Set("3", "The only true wisdom is in knowing you know nothing. - Socrates")
	quotesdb.Set("4", "The greatest wealth is to live content with little. - Plato")
	quotesdb.Set("5", "In the midst of chaos, there is also opportunity. - Sun Tzu")
	quotesdb.Set("6", "Happiness and freedom begin with one principle. Some things are within your control and some are not. - Epictetus")

	resourcesdb := kvstore.NewStore[string, string]()

	return &App{
		commands: command.NewFactory(hc, quotesdb, resourcesdb),
	}
}

func (a *App) Run(ctx context.Context, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	a.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go a.handleConnection(context.Background(), conn)
	}
}

func (a *App) Shutdown(ctx context.Context) {

}

func (a *App) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Minute))
	conn.SetWriteDeadline(time.Now().Add(time.Minute))

	r := textproto.NewReader(bufio.NewReader(conn))
	w := textproto.NewWriter(bufio.NewWriter(conn))

	for {
		quit := false

		select {
		case <-ctx.Done():
			quit = true
		default:
		}

		if quit {
			break
		}

		message, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("closing connection")
				return
			}

			log.Println("failed to read string from connection:", err)
			return
		}

		log.Printf("received a challenge: %q", string(message))

		chunks := strings.Split(message, " ")
		if len(chunks) == 0 {
			continue
		}

		cmd := chunks[0]

		command, err := a.commands.Find(cmd)
		if err != nil {
			log.Println("not found", cmd)
			continue
		}

		result, err := command.Execute(ctx, chunks[1:])
		if err != nil {
			continue
		}

		log.Println("result", result)
		w.PrintfLine(result)
	}
}
