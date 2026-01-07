package main

import (
	"fmt"
	"slices"
)

// https://leetcode.com/problems/contains-duplicate/description/

// HASH SET
// Time Complexity: O(n)
// Space Complexity: O(n)
func containsDuplicateHashing(nums []int) bool {
	seen := make(map[int]struct{})
	for _, num := range nums {
		_, exists := seen[num]
		if exists {
			return true
		}
		seen[num] = struct{}{}
	}
	return false
}

// SORTING
// Time Complexity: O(nlogn)
// Space Complexity: O(1)
func containsDuplicateSorting(nums []int) bool {
	slices.Sort(nums)
	for i := 1; i < len(nums); i++ {
		if nums[i-1] == nums[i] {
			return true
		}
	}
	return false
}

// BRUTE FORCE
// Time Complexity: O(n^2)
// Space Complexity: O(1)
func containsDuplicateBruteForce(nums []int) bool {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] {
				return true
			}
		}
	}

	return false
}

func main() {

	// output true since the element 1 occurs at the indices 0 and 3.
	s1 := []int{1,2,3,1}

	// output false since all elements are distinct.
	s2 := []int{1,2,3,4}

	// output true
	s3 := []int{1,1,1,3,3,4,3,2,4,2}

	fmt.Println("=== CONTAINS DUPLICATES ===")
	fmt.Printf("%v -> %t\n", s1, containsDuplicateBruteForce(s1))
	fmt.Printf("%v -> %t\n", s2, containsDuplicateSorting(s2))
	fmt.Printf("%v -> %t\n", s3, containsDuplicateHashing(s3))

}