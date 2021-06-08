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

func (s *StringSet) ToSlice() []string {
	var result []string
	for k, _ := range s.data {
		result = append(result, k)
	}
	return result
}

func (s *StringSet) Delete(item string) {
	delete(s.data, item)
}
