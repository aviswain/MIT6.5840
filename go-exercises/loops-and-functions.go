package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	guess := 1.0
	for {
		nextGuess := guess - (guess * guess - x) / (2 * guess)

		diff := nextGuess - guess
		if diff < 0 {
			diff = -diff
		}

		if diff < 1e-12 {
			return nextGuess
		}

		guess = nextGuess
	}
}

func main() {
	fmt.Println("My result: ", Sqrt(2))
	fmt.Println("Math library result: ", math.Sqrt(2))
}
