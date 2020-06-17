package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()

	return out
}

func transmitter(in <-chan int) <- chan int {
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()

	return out
}

func merge(done chan struct{}, ch ...<-chan int) <- chan int {
	var wg sync.WaitGroup

	merged := make(chan int)

	waitCh := func(c <- chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case merged <- n:
			case <- done:
				return
			}
		}
	}

	wg.Add(len(ch))

	for _, c := range ch {
		go waitCh(c)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged

}

func main() {
	in := gen(1,2,3,4,5,6)

	done := make(chan struct{}, 2)

	ch1 := transmitter(in)
	ch2 := transmitter(in)

	out := merge(done, ch1, ch2)

	fmt.Println(<- out)

	done <- struct{}{}
	done <- struct{}{}
}