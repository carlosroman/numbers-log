package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"load-test/pkg"
	"math/rand"
	"sync"
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
	errChan := make(chan error, connections)
	ctChan := make(chan struct{}, connections*1000)

	// Count incrementer
	go func() {
		counter := uint64(0)
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ctChan:
				counter++
			case <-ticker.C:
				now := time.Now()
				fmt.Println(fmt.Sprintf("[%v] Current throughput is %v numbers/s using %v connections", now.Format("15:04:05"), counter/10, connections))
				counter = 0
			}
		}
	}()

	for i := 0; i < connections; i++ {
		go func() {
			defer wg.Done()
			client := pkg.NewClient(servAddr)
			err = client.Connect()
			if err != nil {
				errChan <- err
				return
			}
			r := rand.New(rand.NewSource(1337 + int64(i)))
			for {
				n := r.Int31n(1000000000)
				err = client.Send(uint32(n))
				if err != nil {
					errChan <- err
					break
				}
				ctChan <- struct{}{}
			}
		}()
	}
	wg.Wait()

	err = <-errChan
	return err
}
