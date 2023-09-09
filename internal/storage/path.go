package storage

import "strings"

func GetUsersPath() string {
	return "/user"
}

func GetTasksPath() string {
	return "/user"
}

func GetUserPath(id string) string {
	return strings.Join([]string{"/user", id}, "/")
}
