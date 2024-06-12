package main

// Package context defines the Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.

// Incoming requests to a server should create a Context, and outgoing calls to servers should accept a Context. The chain of function calls between them must propagate the Context, optionally replacing it with a derived Context created using WithCancel, WithDeadline, WithTimeout, or WithValue. When a Context is canceled, all Contexts derived from it are also canceled.

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	//ctx := context.Background()

	// Contexts can be used with values, which can be helpfull to trace request ids
	// for example. It also can deal with structs not only strings.
	ctx := context.WithValue(context.Background(), "foo", "bar")
	userID := 10
	val, err := fetchUserData(ctx, userID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result: ", val)
	fmt.Println("took: ", time.Since(start))
}

type Response struct {
	value int
	err   error
}

func fetchUserData(ctx context.Context, userId int) (int, error) {
	// Getting the value of ctx value parameter
	value := ctx.Value("foo")
	fmt.Println(value)
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respch := make(chan Response)

	go func() {
		val, err := fetchThirdPartyStuffSlowly()
		respch <- Response{
			value: val,
			err:   err,
		}

	}()

	for {
		select {
		// It means that the context took more than 200ms to complete.
		// So it timedout.
		// ctx.Done() returns a closed chan <- struct{}
		case <-ctx.Done():
			return 0, fmt.Errorf("Timeout")
		// On this case fetchThirdPartyStuffSlowly() did not take more than 200ms
		// And then it is possible to read values from the respch
		case resp := <-respch:
			return resp.value, resp.err
		}
	}

}

func fetchThirdPartyStuffSlowly() (int, error) {
	time.Sleep(time.Millisecond * 150)

	return 666, nil
}
