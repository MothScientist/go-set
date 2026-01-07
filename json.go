package set

import (
	"encoding/json"
)

// MarshalJSON implements json.Marshaler.
func (s *set[T]) MarshalJSON() ([]byte, error) {
    elements := s.ToSlice()
    return json.Marshal(elements)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *set[T]) UnmarshalJSON(data []byte) error {
	s.mu.Lock()
    defer s.mu.Unlock()

    var elements []T
    if err := json.Unmarshal(data, &elements); err != nil {
        return err
    }

	s.clear()
	s.extendFromSlice(elements)
	return nil
}