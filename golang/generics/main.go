package main

import "fmt"

type Identifiable interface {
	ID() string
}

func IndexByID[T Identifiable](items []T) map[string]T {
	result := make(map[string]T, len(items))
	for _, item := range items {
		result[item.ID()] = item // 约束保证可以调用 ID
	}
	return result
}

type User struct {
	UserID   string
	UserName string
}

func (u User) ID() string { return u.UserID }

func main() {
	users := []User{
		{UserID: "u-1", UserName: "Ada"},
		{UserID: "u-2", UserName: "Rob"},
	}
	index := IndexByID(users)
	fmt.Println(index["u-2"].UserName) // Rob
}
