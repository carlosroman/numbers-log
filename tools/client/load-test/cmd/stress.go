package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"load-test/pkg"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	serverAddress string
	connections   int

	stressCmd = &cobra.Command{
		Use:   "stress [command name]",
		Short: "Runs a stress test at a number-log server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sendNumbers(serverAddress, connections)
		},
	}
)

func init() {
	stressCmd.Flags().
		StringVarP(&serverAddress, "target", "t", "0.0.0.0:4000", "set the server address in form <address>:<port>")
	stressCmd.Flags().
		IntVarP(&connections, "number", "n", 5, "set the number of connections to make")
}

func sendNumbers(servAddr string, connections int) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(connections)
	errChn := make(chan error, connections)

	var counter uint64

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				old := atomic.SwapUint64(&counter, 0)
				fmt.Println(fmt.Sprintf("Current through put is %v numbers/s", old/10))
			}
		}
	}()

	for i := 0; i < connections; i++ {
		go func() {
			defer wg.Done()
			client := pkg.NewClient(servAddr)
			err = client.Connect()
			if err != nil {
				errChn <- err
				return
			}
			r := rand.New(rand.NewSource(1337 + int64(i)))
			for {
				n := r.Int31n(1000000000)
				err = client.Send(uint32(n))
				if err != nil {
					errChn <- err
					break
				}
				atomic.AddUint64(&counter, 1)
			}
		}()
	}
	wg.Wait()

	err = <-errChn
	return err
}
