package paystack

import (
	"context"
	"reflect"
)

// PageFetcher fetches one page of results for a paginated list endpoint,
// given a 1-indexed page number and page size. It returns the decoded
// items for that page along with the call's Response metadata.
type PageFetcher[T any] func(ctx context.Context, page, perPage int) ([]T, *Response, error)

// IteratorOption configures an Iterator.
type IteratorOption[T any] func(*Iterator[T])

// WithMaxPages caps the number of pages an Iterator will fetch, protecting
// callers from accidentally iterating without bound.
func WithMaxPages[T any](n int) IteratorOption[T] {
	return func(it *Iterator[T]) {
		it.maxPages = n
	}
}

// Iterator provides safe iteration over a paginated list endpoint:
//
//	iter := paystack.NewIterator(ctx, 50, client.Transactions.listPage)
//	for iter.Next() {
//	    tx := iter.Item()
//	    // process tx
//	}
//	if err := iter.Err(); err != nil {
//	    // handle
//	}
//
// Iterator stops when the underlying endpoint returns an empty page, when
// it detects the same page being returned twice in a row (guarding against
// a misbehaving API that ignores the page parameter), when a caller-
// supplied page limit is reached, or when the context is canceled.
type Iterator[T any] struct {
	ctx     context.Context
	fetch   PageFetcher[T]
	perPage int
	page    int

	maxPages int

	buffer   []T
	previous []T
	cursor   int

	done bool
	err  error
}

// NewIterator constructs an Iterator over the given PageFetcher. perPage
// values less than or equal to zero default to 50.
func NewIterator[T any](ctx context.Context, perPage int, fetch PageFetcher[T], opts ...IteratorOption[T]) *Iterator[T] {
	if perPage <= 0 {
		perPage = 50
	}
	it := &Iterator[T]{ctx: ctx, fetch: fetch, perPage: perPage, page: 1, cursor: -1}
	for _, opt := range opts {
		opt(it)
	}
	return it
}

// Next advances the iterator to the next item, fetching additional pages
// as needed. It returns false when iteration is complete or an error
// occurred; call Err to distinguish the two.
func (it *Iterator[T]) Next() bool {
	if it.err != nil || it.done {
		return false
	}

	it.cursor++
	if it.cursor < len(it.buffer) {
		return true
	}

	if err := it.ctx.Err(); err != nil {
		it.err = err
		return false
	}

	if it.maxPages > 0 && it.page > it.maxPages {
		it.done = true
		return false
	}

	items, _, err := it.fetch(it.ctx, it.page, it.perPage)
	if err != nil {
		it.err = err
		return false
	}

	if len(items) == 0 {
		it.done = true
		return false
	}

	if it.previous != nil && reflect.DeepEqual(items, it.previous) {
		it.done = true
		return false
	}

	it.buffer = items
	it.previous = items
	it.cursor = 0
	it.page++

	return true
}

// Item returns the current item. It is only valid to call after a call to
// Next that returned true.
func (it *Iterator[T]) Item() T {
	var zero T
	if it.cursor < 0 || it.cursor >= len(it.buffer) {
		return zero
	}
	return it.buffer[it.cursor]
}

// Err returns the first error encountered during iteration, if any.
func (it *Iterator[T]) Err() error {
	return it.err
}
