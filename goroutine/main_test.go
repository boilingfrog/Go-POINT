package main

import (
	"fmt"
	"sync"
	"testing"

	"go.uber.org/goleak"
)

func TestLeakWithGoleak(t *testing.T) {
	defer goleak.VerifyNone(t)
	Goleak()
}

func Goleak() {
	errChannel := make(chan map[string]interface{}, 1)
	normalChannel := make(chan map[string]int, 1)
	endChannel := make(chan bool)

	go func() {
		defer func() {
			endChannel <- true
		}()
		_, _ = listenChan(normalChannel, errChannel, nil, nil)
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(j int) {
			var (
				errMap = map[string]interface{}{}
				norMap = map[string]int{
					"gold": i + 1,
				}
			)
			defer func() {
				errChannel <- errMap
				normalChannel <- norMap
				wg.Done()
			}()
		}(i)
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(j int) {
			var (
				errMap = map[string]interface{}{}
				norMap = map[string]int{
					"gold": 0,
				}
			)
			defer func() {
				errChannel <- errMap
				normalChannel <- norMap
				wg.Done()
			}()
		}(i)
	}
	wg.Wait()
	close(errChannel)
	close(normalChannel)
	<-endChannel
	close(endChannel)
}

func listenChan(normalChannel chan map[string]int, errChannel chan map[string]interface{}, arrayError []map[string]interface{}, arrayNormal []map[string]interface{}) ([]map[string]interface{}, []map[string]interface{}) {
	var (
		errFlag   bool
		norFlag   bool
		goldCount int
	)
	for {
		select {
		case errMap, errFlag := <-errChannel:
			for k, v := range errMap {
				mapError := make(map[string]interface{})
				mapError["type"] = k
				mapError["error"] = v
				if len(mapError) != 0 {
					arrayError = append(arrayError, mapError)
				}
			}
			if !errFlag {
				errChannel = nil
			}
		case norMap, norFlag := <-normalChannel:
			for _, item := range norMap {
				mapNormal := make(map[string]interface{})

				mapNormal["gold"] = item
				if len(norMap) != 0 {
					arrayNormal = append(arrayNormal, mapNormal)
				}
				// 总金币数
				goldCount += item
			}
			if !norFlag {
				normalChannel = nil
			}
		}
		if normalChannel == nil && errChannel == nil {
			goto ForEnd
		}
	}
ForEnd:
	if !errFlag && !norFlag {
		if goldCount >= 0 {
			fmt.Println(goldCount)
		}
	}
	return arrayNormal, arrayError
}
