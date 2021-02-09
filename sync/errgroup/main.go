package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	fmt.Println(ctx)
	eg, _ := errgroup.WithContext(ctx)

	var egg errgroup.Group

	fmt.Println(ctx)
	fmt.Println(eg)
	fmt.Println(egg)

	eg.Go(func() error {
		fmt.Println("11111")
		return errors.New("test1")
	})

	eg.Go(func() error {
		fmt.Println("22222")
		return errors.New("test2")
	})

	fmt.Println(eg.Wait())
}
