package main

import "fmt"

func Sum(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}

	return sum
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	sum := Sum(nums)
	fmt.Println(sum)
}
