package core

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/xpladev/xpla.go/types"
)

var PageRequest *query.PageRequest

// Set default pagination.
func DefaultPagination() *query.PageRequest {
	return &query.PageRequest{
		Key:        []byte(""),
		Offset:     0,
		Limit:      0,
		CountTotal: false,
		Reverse:    false,
	}
}

// Read pagination in the xpla client option
func ReadPageRequest(pagination types.Pagination) (*query.PageRequest, error) {
	pageKey := pagination.PageKey
	offset := pagination.Offset
	limit := pagination.Limit
	countTotal := pagination.CountTotal
	page := pagination.Page
	reverse := pagination.Reverse

	if page > 1 && offset > 0 {
		return nil, types.ErrWrap(types.ErrInvalidRequest, "page and offset cannot be used together")
	}

	if page > 1 {
		offset = (page - 1) * limit
	}

	return &query.PageRequest{
		Key:        []byte(pageKey),
		Offset:     offset,
		Limit:      limit,
		CountTotal: countTotal,
		Reverse:    reverse,
	}, nil
}
