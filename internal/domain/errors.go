package domain
import "errors"

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrDatabase       = errors.New("database error")
	ErrPasswordHash   = errors.New("password hash failed")
)