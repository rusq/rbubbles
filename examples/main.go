package main

import (
	"log"
	"log/slog"
	"os"
)

func main() {
	f, err := os.Create("debug.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	slog.SetDefault(slog.New(
		slog.NewTextHandler(f, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}),
	))
	// wiztest()
	filebrowser()
	// customiseTest()
}
