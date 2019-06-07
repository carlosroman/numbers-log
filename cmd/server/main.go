package main

import (
	"fmt"
	"github.com/carlosroman/numbers-log/internal/pkg/server"
	"os"
)

func main() {

	rec := server.NewRecorder()
	nc := server.NewNumberChecker(rec)
	wr := server.GetWriter("numbers.log")
	defer func() {
		if err := wr.Sync(); err != nil {
			os.Exit(3)
		}
	}()
	h := server.NewHandler(nc, wr)
	s := server.NewServer(5, "localhost", 4000, h)
	if err := s.Start(); err != nil {
		os.Exit(2)
	}
	if err := s.Process(); err != nil {
		os.Exit(1)
	}
	fmt.Println("Done")
}
