package main

import (
	"fmt"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	
	guess := 1.0
	for {
		nextGuess := guess - (guess * guess - x) / (2 * guess)

		diff := nextGuess - guess
		if diff < 0 {
			diff = -diff
		}

		if diff < 1e-12 {
			return nextGuess, nil
		}

		guess = nextGuess
	}
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
