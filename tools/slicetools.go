package tools

import (
	"slices"
	"strings"
)

func TrimSlice(slice []string) []string {
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
	return slice
}

// checks if both slices are of equal length and both contain the same elements, regardless of the order.
func HasSameElements[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for _, s := range s1 {
		if !slices.Contains(s2, s) {
			return false
		}
	}
	return true
}

func AppendUnique[T comparable](to []T, s ...T) []T {
	target := EnsureUnique(to)
	for _, ss := range s {
		if slices.Contains(target, ss) {
			continue
		}
		target = append(target, ss)
	}
	return target
}

func EnsureUnique[T comparable](s1 []T) []T {
	ss := make([]T, 0, len(s1))
	for _, e := range s1 {
		if slices.Contains(ss, e) {
			continue
		}
		ss = append(ss, e)
	}
	return ss
}

func NotIn[T comparable](s1, notIn []T) []T {
	ss := make([]T, 0, len(s1))
	for _, e := range s1 {
		if slices.Contains(notIn, e) {
			continue
		}
		ss = append(ss, e)
	}
	return ss
}
