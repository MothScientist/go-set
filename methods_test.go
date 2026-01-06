package set

import (
	"testing"
)

func TestCloneTableDriven(t *testing.T) {
	tests := []struct {
		testName  string
		inputSet  *set[int]
	}{
		{"{}", NewEmptySet[int](0)},
		{"{1}", NewSet(1)},
		{"{1,2,3}", NewSet(1, 2, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			newSet := Clone(tt.inputSet)
			if !Equal(newSet, tt.inputSet) {
				t.Errorf("got %s, want %s", newSet, tt.inputSet)
			}
		})
	}
}

func TestDifferenceTableDriven(t *testing.T) {
	tests := []struct {
		testName  string
		originalSet  *set[int]
		secondSet *set[int]
		resultSet *set[int]
	}{
		{"{}-{}->{}", NewEmptySet[int](0), NewEmptySet[int](0), NewEmptySet[int](0)},
		{"{}-{1}->{}", NewEmptySet[int](0), NewSet(1), NewEmptySet[int](0)},
		{"{1}-{}->{1}", NewSet(1), NewEmptySet[int](0), NewSet(1)},
		{"{1,2,3}-{4,5,6}->{1,2,3}", NewSet(1, 2, 3), NewSet(4, 5, 6), NewSet(1, 2, 3)},
		{"{1,2,3}-{1,2,3}->{}", NewSet(1, 2, 3), NewSet(1, 2, 3), NewEmptySet[int](0)},
		{"{1,2,3}-{2,3}->{1}", NewSet(1, 2, 3), NewSet(2, 3), NewSet(1)},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			newSet := Difference(tt.originalSet, tt.secondSet)
			if !Equal(newSet, tt.resultSet) {
				t.Errorf("got %s, want %s", newSet, tt.resultSet)
			}
		})
	}
}

func TestIntersectionTableDriven(t *testing.T) {
	tests := []struct {
		testName  string
		originalSet  *set[int]
		secondSet *set[int]
		resultSet *set[int]
	}{
		{"{}&{}->{}", NewEmptySet[int](0), NewEmptySet[int](0), NewEmptySet[int](0)},
		{"{}&{1}->{}", NewEmptySet[int](0), NewSet(1), NewEmptySet[int](0)},
		{"{1}&{}->{}", NewSet(1), NewEmptySet[int](0), NewEmptySet[int](0)},
		{"{1,2,3}&{4,5,6}->{}", NewSet(1, 2, 3), NewSet(4, 5, 6), NewEmptySet[int](0)},
		{"{1,2,3}&{1,2,3}->{1,2,3}", NewSet(1, 2, 3), NewSet(1, 2, 3), NewSet(1, 2, 3)},
		{"{1,2,3}&{2,3}->{2,3}", NewSet(1, 2, 3), NewSet(2, 3), NewSet(2,3)},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			newSet := Intersection(tt.originalSet, tt.secondSet)
			if !Equal(newSet, tt.resultSet) {
				t.Errorf("got %s, want %s", newSet, tt.resultSet)
			}
		})
	}
}

func TestUnionTableDriven(t *testing.T) {
	tests := []struct {
		testName  string
		originalSet  *set[int]
		secondSet *set[int]
		resultSet *set[int]
	}{
		{"{}|{}->{}", NewEmptySet[int](0), NewEmptySet[int](0), NewEmptySet[int](0)},
		{"{}|{1}->{1}", NewEmptySet[int](0), NewSet(1), NewSet(1)},
		{"{1}|{}->{1}", NewSet(1), NewEmptySet[int](0), NewSet(1)},
		{"{1,2,3}|{4,5,6}->{1,2,3,4,5,6}", NewSet(1, 2, 3), NewSet(4, 5, 6), NewSet(1, 2, 3, 4, 5, 6)},
		{"{1,2,3}|{1,2,3}->{1,2,3}", NewSet(1, 2, 3), NewSet(1, 2, 3), NewSet(1, 2, 3)},
		{"{1,2,3}|{2,3}->{1,2,3}", NewSet(1, 2, 3), NewSet(2, 3), NewSet(1, 2,3)},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			newSet := Union(tt.originalSet, tt.secondSet)
			if !Equal(newSet, tt.resultSet) {
				t.Errorf("got %s, want %s", newSet, tt.resultSet)
			}
		})
	}
}

func TestSymmetricDifferenceTableDriven(t *testing.T) {
	tests := []struct {
		testName  string
		originalSet  *set[int]
		secondSet *set[int]
		resultSet *set[int]
	}{
		{"{}^{}->{}", NewEmptySet[int](0), NewEmptySet[int](0), NewEmptySet[int](0)},
		{"{}^{1}->{1}", NewEmptySet[int](0), NewSet(1), NewSet(1)},
		{"{1}^{}->{1}", NewSet(1), NewEmptySet[int](0), NewSet(1)},
		{"{1,2,3}^{4,5,6}->{1,2,3,4,5,6}", NewSet(1, 2, 3), NewSet(4, 5, 6), NewSet(1, 2, 3, 4, 5, 6)},
		{"{1,2,3}^{1,2,3}->{}", NewSet(1, 2, 3), NewSet(1, 2, 3), NewEmptySet[int](0)},
		{"{1,2,3}^{2,3}->{1}", NewSet(1, 2, 3), NewSet(2, 3), NewSet(1)},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			newSet := SymmetricDifference(tt.originalSet, tt.secondSet)
			if !Equal(newSet, tt.resultSet) {
				t.Errorf("got %s, want %s", newSet, tt.resultSet)
			}
		})
	}
}