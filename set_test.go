package set

import (
	"sync"
	"testing"
)

func TestInitEmptySet(t *testing.T) {
	set := NewEmptySet[int](0)
	lenSet := len(set.s)
	if lenSet != 0 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", lenSet, 0)
	}
}

func TestInitSet(t *testing.T) {
	set := NewSet(1, 2, 3)
	lenSet := len(set.s)
	if lenSet != 3 {
		t.Errorf("Result was incorrect, got: %d, want: %d.", lenSet, 3)
	}
}

func TestStringInterface(t *testing.T) {
	tests := []struct {
		testName  string
		inputSet  *set[string]
		lenString int
	}{
		{"{}", NewEmptySet[string](0), 5},      // Set[]
		{"{1}", NewSet("1"), 6},                // Set[1]
		{"{1,2,3}", NewSet("1", "2", "3"), 10}, // Set[1 2 3] (or Set[2 3 1] etc...)
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			l := len(tt.inputSet.String())
			if l != tt.lenString {
				t.Errorf("got %d, want %d", l, tt.lenString)
			}
		})
	}
}

func TestEmptyTableDriven(t *testing.T) {
	tests := []struct {
		testName string
		inputSet *set[int]
		isEmpty  bool
	}{
		{"[]make(0)", NewEmptySet[int](0), true},
		{"[]make(5)", NewEmptySet[int](5), true},
		{"[1]", NewSet(1), false},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			setIsEmpty := tt.inputSet.Empty()
			if setIsEmpty != tt.isEmpty {
				t.Errorf("got %t, want %t", setIsEmpty, tt.isEmpty)
			}
		})
	}
}

func TestSizeTableDriven(t *testing.T) {
	tests := []struct {
		testName string
		inputSet *set[int]
		sizeSet  int
	}{
		{"[]make(0)", NewEmptySet[int](0), 0},
		{"[]make(5)", NewEmptySet[int](5), 0},
		{"[1]", NewSet(1), 1},
		{"[1,2,3]", NewSet(1, 2, 3), 3},
		{"[1(1,1,1)]", NewSet(1, 1, 1, 1), 1},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			lenSet := tt.inputSet.Size()
			if lenSet != tt.sizeSet {
				t.Errorf("got %d, want %d", lenSet, tt.sizeSet)
			}
		})
	}
}

func TestInsertTableDriven(t *testing.T) {
	tests := []struct {
		testName    string
		inputSet    *set[int]
		insertValue int
		wantSet     *set[int]
	}{
		{"1->[1,2,3]", NewSet(1, 2, 3), 1, NewSet(1, 2, 3)},
		{"4->[1,2,3]", NewSet(1, 2, 3), 4, NewSet(1, 2, 3, 4)},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.inputSet.Insert(tt.insertValue)
			if !Equal(tt.inputSet, tt.wantSet) {
				t.Errorf("got %s, want %s", tt.inputSet, tt.wantSet)
			}
		})
	}
}

func TestExistsTableDriven(t *testing.T) {
	tests := []struct {
		testName      string
		inputSet      *set[int]
		valueToSearch int
		want          bool
	}{
		{"[]?1-make(0)", NewEmptySet[int](0), 1, false},
		{"[]?0-make(5)", NewEmptySet[int](0), 0, false},
		{"[1,2,3]?1", NewSet(1, 2, 3), 1, true},
		{"[1,2,3]?5", NewSet(1, 2, 3), 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := tt.inputSet.Exists(tt.valueToSearch)
			if got != tt.want {
				t.Errorf("got %t, want %t", got, tt.want)
			}
		})
	}
}

func TestDeleteTableDriven(t *testing.T) {
	tests := []struct {
		testName    string
		inputSet    *set[int]
		deleteValue int
		wantSet     *set[int]
	}{
		{"[]-1->[]", NewEmptySet[int](0), 1, NewEmptySet[int](0)},
		{"[1,2,3]-3->[1,2]", NewSet(1, 2, 3), 3, NewSet(1, 2)},
		{"[1,2,3]-5->[1,2,3]", NewSet(1, 2, 3), 5, NewSet(1, 2, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.inputSet.Delete(tt.deleteValue)
			if !Equal(tt.inputSet, tt.wantSet) {
				t.Errorf("got %s, want %s", tt.inputSet, tt.wantSet)
			}
		})
	}
}

func TestClearTableDriven(t *testing.T) {
	tests := []struct {
		testName string
		inputSet *set[int]
	}{
		{"{}-make(0)", NewEmptySet[int](0)},
		{"{}-make(3)", NewEmptySet[int](3)},
		{"{1}", NewSet(1)},
		{"{1,2,3}", NewSet(1, 2, 3)},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.inputSet.Clear()
			if tt.inputSet.Size() != 0 {
				t.Errorf("got [], want %s", tt.inputSet)
			}
		})
	}
}

func TestExtendFromSliceIntTableDriven(t *testing.T) {
	tests := []struct {
		testName    string
		inputSet    *set[int]
		extendSlice []int
		wantSet     *set[int]
	}{
		{"{}+[1,2,3]->{1,2,3}", NewEmptySet[int](3), []int{1, 2, 3}, NewSet(1, 2, 3)},
		{"{}+[]->{}", NewEmptySet[int](0), []int{}, NewEmptySet[int](0)},
		{"{1}+[]->{1}", NewSet(1), []int{}, NewSet(1)},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.inputSet.ExtendFromSlice(tt.extendSlice)
			if !Equal(tt.inputSet, tt.wantSet) {
				t.Errorf("got %s, want %s", tt.inputSet, tt.wantSet)
			}
		})
	}
}

func TestExtendFromMapIntTableDriven(t *testing.T) {
	tests := []struct {
		testName  string
		inputSet  *set[int]
		extendMap map[int]any
		wantSet   *set[int]
	}{
		{"{}+[1]->{1}", NewEmptySet[int](0), map[int]any{1: "123"}, NewSet(1)},
		{"{1}+[1,2,3]->{1,2,3}", NewSet(1), map[int]any{1: "123", 2: true, 3: 4}, NewSet(1, 2, 3)},
		{"{1}+[]->{1}", NewSet(1), map[int]any{}, NewSet(1)},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.inputSet.ExtendFromMap(tt.extendMap)
			if !Equal(tt.inputSet, tt.wantSet) {
				t.Errorf("got %s, want %s", tt.inputSet, tt.wantSet)
			}
		})
	}
}

func TestDataRace(t *testing.T) {
	var wg sync.WaitGroup
	count := 100
	set := NewEmptySet[int](count)
	for i := range count {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			set.Insert(i)
		}(i)
	}
	for i := range count {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			set.Delete(i)
		}(i)
	}
	wg.Wait()
}
