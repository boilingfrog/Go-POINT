package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	eg := errgroup.Group{}
	eg.Go(func() (err error) {
		// 处理业务
		err = test1()
		return err
	})

	eg.Go(func() (err error) {
		// 处理业务
		err = test1()
		return err
	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}
}

func test1() error {
	return errors.New("test2")
}

func ttt() int {
	a := 1
	if a > 1 {
		return a
	} else {
		return 2
	}
}

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result string
type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
	return func(_ context.Context, query string) (Result, error) {
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

func ExampleGroup_parallel() {
	Google := func(ctx context.Context, query string) ([]Result, error) {
		g, ctx := errgroup.WithContext(ctx)

		searches := []Search{Web, Image, Video}
		results := make([]Result, len(searches))
		for i, search := range searches {
			i, search := i, search
			g.Go(func() error {
				result, err := search(ctx, query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
