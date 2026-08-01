package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"
)

func main() {
	alloc := make([]int, 10)

	alloc = append(alloc, int(1))

	done := make(chan bool)

	a := unsafe.Sizeof(alloc)
	fmt.Println("Current size of slice is ", a)

	tm := time.NewTicker(time.Second)
	defer tm.Stop()

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	wg.Add(12)

	go func(){
		defer close(done)
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Printf("Application stopping")
				return
			case <-tm.C:
				fmt.Println("Current size of application, ", unsafe.Sizeof(alloc) - uintptr(len(alloc)) * unsafe.Sizeof(int(0)))
			}
		}
	}()

	go func(){
		defer wg.Done()
		time.Sleep(10 * time.Second)
		done <- true
	}()

	for range 10{
		mu.Lock()
		go addElement(&alloc, wg)
		mu.Unlock()

		time.Sleep(100 * time.Millisecond)
	}

	wg.Wait()
}

func addElement(alloc *[]int, wg *sync.WaitGroup) {
	defer wg.Done()
	*alloc = append(*alloc, int(1))
}