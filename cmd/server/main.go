package main

import (
	"context"
	"os"

	"github.com/maksadbek/wordofwisdom/cmd/server/app"
)

var (
	addr = os.Getenv("ADDR")
)

func main() {
	app := app.NewApp()
	app.Run(context.Background(), addr)

	defer app.Shutdown(context.Background())
}
