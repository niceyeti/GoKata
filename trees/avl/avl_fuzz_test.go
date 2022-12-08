package avl

import (
	"testing"
)

// TODO: I intentionally left this as failing, only fooling
// around with fuzzing for now. Fuzzing seems useful for data structures
// solely for exercising large volumes of random api calls in order to discover
// edge cases that a developer/tester did not consider.
func FuzzInsertion(f *testing.F) {
	modulus := 10067
	// Need to track the used ints, since the tree disallows dupes.
	usedInts := make(map[int]bool, modulus)

	testcases := []int{1, 2, 3, 4, 5}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
		usedInts[tc] = true
	}

	tr := NewTree()
	f.Fuzz(func(t *testing.T, in int) {
		n := in % modulus
		if _, ok := usedInts[n]; !ok {
			usedInts[n] = true
			err := tr.Insert(in)
			t.Logf("Input %d", in)
			if err != nil {
				t.Errorf("Inserted %d but got err %v", in, err)
			}
		}
	})
}
