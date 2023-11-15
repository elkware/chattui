package internal

import "os/user"

func GetUsernameFromOS() string {
	u, _ := user.Current()
	return u.Username
}
