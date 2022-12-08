package avl

import (
	"fmt"
	"testing"
)

// TODO: I intentionally left this as failing. Only fooling
// around with fuzzing for now. It seems useful for data structures
// solely for exercising lots and lots of api calls in order to catch
// edge cases a developer/tester that a developer did not think about.
func FuzzInsertion(f *testing.F) {
	testcases := []int{1, 2, 3, 4, 5}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	tr := NewTree()
	f.Fuzz(func(t *testing.T, in int) {
		err := tr.Insert(in)
		fmt.Println(in)
		t.Logf("Input %d", in)
		if err != nil {
			t.Errorf("Inserted %d but got err %v", in, err)
		}
	})
}
