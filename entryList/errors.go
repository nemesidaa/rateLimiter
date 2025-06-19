package entrylist

import "errors"

var (
	ErrDownStated       = errors.New("Luckily you caught STW! Your request `ll be handled in a bit...")
	ErrPersonalOverflow = errors.New("You`re reached a limit of requests this time. Try again a bit later")
)
