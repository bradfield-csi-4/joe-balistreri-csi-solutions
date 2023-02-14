package node

import (
	"fmt"
	"testing"
	"time"
)

func TestScan(t *testing.T) {
	t.Run("benchmark for scanning ratings", func(t *testing.T) {
		tablenameToFilename = Row{
			"movies":  "../data/ml-20m/movies.csv",
			"ratings": "../data/ml-20m/ratings.csv",
		}
		n := NewScanNode("ratings")

		s := time.Now()
		for n.Next() != nil {
		}
		dur := time.Since(s)
		fmt.Printf("It took %s to read the ratings table\n", dur)
	})
}
