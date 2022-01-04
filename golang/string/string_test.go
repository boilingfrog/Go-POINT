package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func String() string {
	var s string
	s += "hello" + "\n"
	s += "world" + "\n"
	s += "今天的天气很不错的"
	return s
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		String()
	}
}

func StringFmt() string {
	return fmt.Sprintf("%s %s %s", "hello", "world", "今天的天气很不错的")
}

func BenchmarkStringFmt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringFmt()
	}
}

func StringJoin() string {
	var s []string
	s = append(s, "hello ")
	s = append(s, "world ")
	s = append(s, "今天的天气很不错的 ")

	return strings.Join(s, "")
}

func BenchmarkStringJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringJoin()
	}
}

func StringBuffer() string {
	var s bytes.Buffer
	s.WriteString("hello ")
	s.WriteString("world ")
	s.WriteString("今天的天气很不错的 ")

	return s.String()
}

func BenchmarkStringBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringBuffer()
	}
}

func StringBuilder() string {
	var s strings.Builder
	s.WriteString("hello ")
	s.WriteString("world ")
	s.WriteString("今天的天气很不错的 ")

	return s.String()
}

func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringBuilder()
	}
}
