package main

import (
	"fmt"
	"github.com/carlosroman/numbers-log/internal/pkg/server"
	"os"
)

func main() {

	r := server.NewNumberChecker()
	l := server.GetWriter("numbers.log")
	defer func() {
		if err := l.Sync(); err != nil {
			os.Exit(3)
		}
	}()
	h := server.NewHandler(r, l)
	s := server.NewServer(5, "localhost", 4000, h)
	if err := s.Start(); err != nil {
		os.Exit(2)
	}
	if err := s.Process(); err != nil {
		os.Exit(1)
	}
	fmt.Println("Done")
}
