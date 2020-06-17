package main

import (
	"fmt"
	"sync"
)

func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums))

	for _, n := range nums {
		out <- n
	}
	close(out)

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

func merge(ch ...<-chan int) <- chan int {
	var wg sync.WaitGroup

	merged := make(chan int, 1)

	waitCh := func(c <- chan int) {
		for n := range c {
			merged <- n
		}

		wg.Done()
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
	in2 := gen(10,20,30)

	ch1 := transmitter(in)
	ch2 := transmitter(in2)

	for n := range merge(ch1, ch2) {
		fmt.Println(n)
	}
}