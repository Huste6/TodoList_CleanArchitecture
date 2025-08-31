package common

import "fmt"

const (
	CurrentUser = "current_user"
)

func Recover() {
	if r := recover(); r != nil {
		fmt.Println("Recovered: ", r)
	}
}

type TokenPayLoad struct {
	UID   int    `json:"user_id"`
	URole string `json:"role"`
}

func (p TokenPayLoad) UserId() int {
	return p.UID
}

func (p TokenPayLoad) Role() string {
	return p.URole
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}
