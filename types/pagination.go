package types

type Pagination struct {
	PageKey    string
	Offset     uint64
	Limit      uint64
	CountTotal bool
	Page       uint64
	Reverse    bool
}
