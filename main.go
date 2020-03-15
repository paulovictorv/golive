package main

import (
	"fmt"
	"goclip.com.br/golive/infrastructure"
	"sync"
)

func main() {
	var waitGroup sync.WaitGroup
	progress := make(chan string)
	complete := make(chan int)

	waitGroup.Add(1)

	go infrastructure.CreateEnv("env.Bucket", "env.Domain", progress, complete)
	go func() {
		for {
			select {
			case val := <-progress:
				fmt.Printf(val)
			case <-complete:
				waitGroup.Done()
				break
			}
		}
	}()

	waitGroup.Wait()
}
