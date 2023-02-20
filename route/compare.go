package route

import (
	"fmt"
)

type Compare interface {
	Key() string
	Compare(o Compare) []string
}

func Diff(a1 []Compare, a2 []Compare) []string {
	m1 := make(map[string]Compare, len(a1))
	for _, a := range a1 {
		m1[a.Key()] = a
	}
	m2 := make(map[string]Compare, len(a2))
	for _, a := range a2 {
		m2[a.Key()] = a
	}
	diff := make([]string, 0)
	for k2, v2 := range m2 {
		if v1 := m1[k2]; v1 != nil {
			diff = append(diff, v1.Compare(v2)...)
			delete(m1, k2)
			delete(m2, k2)
		}
	}
	add := make([]string, 0)
	for k2 := range m2 {
		add = append(add, fmt.Sprintf("Add:%s", k2))
	}
	del := make([]string, 0)
	for k1 := range m1 {
		del = append(del, fmt.Sprintf("Delete:%s", k1))
	}
	return append(add, append(del, diff...)...)
}
