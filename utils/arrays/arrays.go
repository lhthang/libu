package arrays

import (
	"fmt"
	"reflect"
	"sort"
)

// Distinct has complexity: O(n^2)
func Distinct(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		if Contains(b, el) {
			set = append(set, el)
		}
	}

	return set
}

// Sorted has complexity: O(n * log(n)), a needs to be sorted
func Sorted(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		idx := sort.Search(bv.Len(), func(i int) bool {
			return bv.Index(i).Interface() == el
		})
		if idx < bv.Len() && bv.Index(idx).Interface() == el {
			set = append(set, el)
		}
	}

	return set
}

// Hash has complexity: O(n * x) where x is a factor of hash function efficiency (between 1 and 2)
func Hash(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	hash := make(map[interface{}]bool)
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		hash[el] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).Interface()
		if _, found := hash[el]; found {
			set = append(set, el)
		}
	}

	return set
}

func Contains(a interface{}, e interface{}) bool {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true
		}
	}
	return false
}

func Intersection(a, b []string) (c []string) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func Union(a, b []string) (c []string) {
	a = append(a, b...)

	seen := make(map[string]struct{}, len(a))
	j := 0
	for _, s := range a {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		a[j] = s
		j++
	}

	mergedSlice := a[:j]
	return mergedSlice
}

// difference returns the elements in `a` that aren't in `b`.
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func Keys(m interface{}) (keys []interface{}) {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		fmt.Errorf("input type not a map: %v", v)
	}

	for _, k := range v.MapKeys() {
		keys = append(keys, k.Interface())
	}
	return keys

}

