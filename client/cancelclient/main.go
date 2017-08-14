package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

func work(ctx context.Context) error {
	defer wg.Done()
	c := make(chan struct {
		r   *http.Response
		err error
	}, 1)

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest("GET", "http://localhost:39090/api/hello", nil)
	go func() {
		resp, err := client.Do(req)
		pack := struct {
			r   *http.Response
			err error
		}{resp, err}
		c <- pack
	}()

	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c
		fmt.Println("Client Canceled!")
		return ctx.Err()
	case ret := <-c:
		err := ret.err
		resp := ret.r
		if err != nil {
			fmt.Println("Error", err)
			return err
		}
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Server said: %s\n", out)
	}
	return nil
}

func retrive() error {
	fmt.Println("calling http rest service!")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	wg.Add(1)
	go work(ctx)
	wg.Wait()
	fmt.Println("Finished calling!")
	return ctx.Err()
}

func main() {
	fmt.Println("*** Start Calling the server!")
	for {
		err := retrive()
		if err != nil {
			fmt.Println("The end!")
			return
		}
	}
}
