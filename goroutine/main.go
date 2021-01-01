package main

import (
	"fmt"
	"time"
)

func main() {
	// case1()
	case6()
}

func channelNoReceiver() {
	ch := make(chan int)
	go func() {
		ch <- 1
		fmt.Println(111)
	}()
}

func channelNoProducer() {
	ch := make(chan int, 1)
	go func() {
		<-ch
		fmt.Println(111)
	}()
}

func case6() {
	names := []string{"小白", "小明", "小红", "小张"}
	for _, name := range names {
		go func(who string) {
			fmt.Println("名字", who)
		}(name)
	}
	time.Sleep(time.Millisecond)
}

func case5() {
	names := []string{"小白", "小明", "小红", "小张"}
	for _, name := range names {
		go func() {
			fmt.Println("名字", name)
		}()
	}
	time.Sleep(time.Millisecond)
}

func case4() {
	name := "小白"
	go func() {
		fmt.Println(name)
	}()
	time.Sleep(time.Millisecond)
	name = "小李"
}

func case3() {
	name := "小白"
	go func() {
		fmt.Println(name)
	}()
	name = "小李"
	time.Sleep(time.Millisecond)

}

func case2() {
	go func() {
		fmt.Println("goroutine", 1)
	}()
	go func() {
		fmt.Println("goroutine", 2)
	}()
	go func() {
		fmt.Println("goroutine", 3)
	}()
}

func case1() {
	names := []string{"小白", "小李", "小张"}
	for _, name := range names {
		go func() {
			fmt.Println(name)
		}()
		time.Sleep(time.Second)
	}
}
