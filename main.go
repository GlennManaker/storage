package main

import (
	"fmt"
	"my_storage/pkg/storage"
	"time"
)

func main() {
	storage := storage.NewStorage[int, int]()
	storage.Start()

	now := time.Now()

	fmt.Println(now)

	for t := range 20000000 {
		storage.Set(t, -t, now.Add(time.Second*30))
		//fmt.Println(storage.Get(key))
	}

	fmt.Println(time.Now().Sub(now))

	time.Sleep(30 * time.Second)

	fmt.Println(storage.GetRoot().GetRootValue())
}
