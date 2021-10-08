package entities

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"name"`
	Password string `json:"price"`
	Age      int    `json:"quantity"`
	Male     bool   `json:"male"`
}

type Login struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

func (l Login) GetUsername() string {
	return l.Username
}
func (l *Login) UpdateUsername(username string) {
	l.Username = username

}
