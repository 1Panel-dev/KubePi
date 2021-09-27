package token

import "time"

//TtyParameter kubectl tty param
type TtyParameter struct {
	Title string
	Arg   string
}

//interface that defines token cache behavior
type Cache interface {
	Get(token string) *TtyParameter
	Delete(token string) error
	Add(token string, param *TtyParameter, d time.Duration) error
}
