package main

import "fmt"

type Student struct {
	Name  string
	Score int
}

func (s *Student) PrintScore() {
	fmt.Println("学生", s.Name, "年纪", s.Score)
}

func NewStudent(name string, score int) *Student {
	return &Student{
		Name:  name,
		Score: score,
	}
}

func main() {
	std := NewStudent("小明", 12)
	std.PrintScore()
}
