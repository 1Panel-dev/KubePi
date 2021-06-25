package collectons

type StringSet struct {
	data map[string]interface{}
}

func NewStringSet() *StringSet {
	ss := StringSet{}
	ss.data = map[string]interface{}{}
	return &ss
}

func (s *StringSet) Add(item string) {
	s.data[item] = nil
}

func (s *StringSet) Exists(item string) bool {
	_, ok := s.data[item]
	return ok
}

func (s *StringSet) Difference(items []string) []string {
	result := make([]string, 0)
	for k, _ := range s.data {
		exists := false
		for i := range items {
			if k == items[i] {
				exists = true
			}
		}
		if !exists {
			result = append(result, k)
		}
	}
	return result
}

func (s *StringSet) ToSlice() []string {
	result := make([]string, 0)
	for k, _ := range s.data {
		result = append(result, k)
	}
	return result
}

func (s *StringSet) Delete(item string) {
	delete(s.data, item)
}
