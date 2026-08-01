package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

func main() {
	alloc := make([]int, 1, 15)
	alloc = append(alloc, int(1))

	done := make(chan bool)

	a := unsafe.Sizeof(alloc) * uintptr(len(alloc))
	fmt.Println("Current size of slice is ", a)

	tm := time.NewTicker(time.Second)
	defer tm.Stop()

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	wg.Add(202 + 10)

	go func(alloc *[]int) {
		defer close(done)
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Printf("Application stopping")
				return
			case <-tm.C:
				fmt.Println("Current size of application, ", unsafe.Sizeof(alloc)*uintptr(len(*alloc)))
			}
		}
	}(&alloc)

	go func() {
		defer wg.Done()
		time.Sleep(15 * time.Second)
		done <- true
	}()

	wg.Go(func() {
		for range 100 {
			mu.Lock()
			go addElement(&alloc, wg)
			mu.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	})


	wg.Go(func() {
		for range 100 {
			mu.Lock()
			go deleteElement(&alloc, wg)
			mu.Unlock()
			time.Sleep(110 * time.Millisecond)
		}
	})

	// booster logic
	wg.Go(func() {
		for range 10 {
			mu.Lock()
			go addElement(&alloc, wg)
			mu.Unlock()
			time.Sleep(45 * time.Millisecond)
		}
	})


	wg.Wait()
}

func addElement(alloc *[]int, wg *sync.WaitGroup) {
	defer wg.Done()
	*alloc = append(*alloc, int(1))
}

func deleteElement(alloc *[]int, wg *sync.WaitGroup){
	defer wg.Done()
	if (len(*alloc) != 0){
		*alloc = (*alloc)[:len(*alloc)-1]
	}
}
