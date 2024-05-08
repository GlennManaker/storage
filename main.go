package main

import (
	"fmt"
	"my_storage/pkg/storage"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	storage := storage.NewStorage[int, int]()
	storage.Start()

	now := time.Now()

	wg.Add(2)

	go func() {
		defer wg.Done()
		for t := range 2000000 {
			storage.Set(t, -t, now.Add(time.Second*10))
			//fmt.Println(storage.Get(key))
		}
	}()

	go func() {
		defer wg.Done()
		for t := range 2000000 {
			storage.Set(t+200000, -t-20000, now.Add(time.Second*10))
			//fmt.Println(storage.Get(key))
		}
	}()

	wg.Wait()

	fmt.Println(time.Now().Sub(now))

	time.Sleep(15 * time.Second)

	fmt.Println(storage.GetRoot())
}
