package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func work(ctx context.Context) error {
	defer wg.Done()
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("I'm working days and nights!")
		case <-ctx.Done():
			fmt.Println("I've got a cancel from client! Ops")
			return ctx.Err()
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	fmt.Println("Hey, Let's do some work, together")
	wg.Add(1)
	go work(ctx)
	wg.Wait()

	fmt.Println("Finished! Enjoy your weekend!")
}
