package entrylist

type EntryListShape interface {
	IncrementFor(entry string) error
	Reset() error
	// Respawn() error
}
