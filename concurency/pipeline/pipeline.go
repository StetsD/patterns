package main

import (
	"fmt"
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

func transmitter(in <- chan int) <-chan int{
	out := make(chan int)

	go func() {
		for n := range in {
			out <- n * 2
		}
		close(out)
	}()

	return out
}

func main() {

	for n := range transmitter(transmitter(gen(1,2,3,4,5,6,7,8,9))) {
		fmt.Println(n)
	}
}
