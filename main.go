package main

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/datastore"
)

const ItemKind = "item"

type Item struct {
	ID    string `datastore:"id"`
	Title string `datastore:"title"`
}

type Store struct {
	client *datastore.Client
	mu     sync.Mutex
}

// PutTx creates or updates a single item in datastore using a transaction.
// This is a simple example just to demonstrate the issue.
func (s *Store) PutTx(ctx context.Context, item *Item, withLock bool) error {
	if withLock {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	tx, err := s.client.NewTransaction(ctx)
	if err != nil {
		return fmt.Errorf("error creating transaction: %v", err)
	}
	k := datastore.NameKey(ItemKind, item.ID, nil)
	_, err = tx.Put(k, item)
	if err != nil {
		return err
	}
	_, err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("error rolling back transaction: %v", rbErr)
		}
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func main() {

}
