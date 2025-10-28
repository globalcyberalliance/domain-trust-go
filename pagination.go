package client

import (
	"context"
)

type (
	// Iterator is a generic pagination iterator.
	Iterator[T any] struct {
		ctx       context.Context
		err       error
		fetchPage PageFetcher[T]
		finished  bool
		index     int
		nextToken string
		page      []T
	}

	// PageFetcher fetches one page of items and returns items, nextPageToken, and an error.
	PageFetcher[T any] func(ctx context.Context, pageToken string) ([]T, string, error)
)

// Err returns any error that occurred during iteration.
func (it *Iterator[T]) Err() error {
	return it.err
}

// Next advances the iterator and loads the next page if needed. It returns true if there is a next value.
func (it *Iterator[T]) Next() bool {
	if it.err != nil || it.finished {
		return false
	}

	// if no page or exhausted, fetch
	if it.page == nil || it.index >= len(it.page) {
		it.page, it.nextToken, it.err = it.fetchPage(it.ctx, it.nextToken)

		if it.err != nil {
			it.finished = true
			return false
		}

		if len(it.page) == 0 {
			it.finished = true
			return false
		}

		it.index = 0
	}

	it.index++

	return true
}

// Value returns the current element.
func (it *Iterator[T]) Value() T {
	return it.page[it.index-1]
}
