package main

import "fmt"

func twoSum(nums []int, target int) []int {
	seen := make(map[int]int)

	for i, value := range nums {
		diff := target - value

		fmt.Println("Diff is: ", diff)

		if index, ok := seen[diff]; ok {
			fmt.Println("found!")
			return []int{index, i}
		} else {
			seen[value] = i
		}

		fmt.Println("Seen is: ", seen)
	}
	return []int{}
}

func main() {
	solution := twoSum([]int{1, 2, 3}, 5)

	fmt.Println("Solution is: ", solution)
}
