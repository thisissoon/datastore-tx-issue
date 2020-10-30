package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"cloud.google.com/go/datastore"

	"github.com/thisissoon/datastore-tx-issue/datastoretest"
)

func initStore(ctx context.Context, t *testing.T) (*Store, *datastore.Client, func()) {
	if !strings.Contains(os.Getenv("DATASTORE_HOST"), "localhost") {
		t.Fatal(fmt.Errorf("Datastore emulator must be configured to run integration tests"))
	}
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	s := &Store{
		client: client,
	}
	return s, client, func() {
		datastoretest.Cleanup(ctx, t, client, ItemKind)
		client.Close()
	}
}

var TestItem = &Item{
	ID:    "d9d6aed8-0623-4f94-9cf6-eefae75a15cf",
	Title: "title1",
}

func TestStore_PutTx(t *testing.T) {
	tests := []struct {
		name     string
		useLock  bool
		putItems []*Item
	}{
		{
			name:    "no lock",
			useLock: false,
			putItems: []*Item{
				TestItem,
				TestItem,
			},
		},
		{
			name:    "with lock",
			useLock: true,
			putItems: []*Item{
				TestItem,
				TestItem,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s, client, cleanup := initStore(ctx, t)
			defer cleanup()
			var wg sync.WaitGroup
			var errors []error
			for i, item := range tt.putItems {
				wg.Add(1)
				go func(index int, item *Item) {
					defer wg.Done()
					err := s.PutTx(ctx, item, tt.useLock)
					if err != nil {
						errors = append(errors, fmt.Errorf("error updating item %d: %v", index, err))
						return
					}
				}(i, item)
			}
			wg.Wait()
			for _, err := range errors {
				t.Fatal(err)
			}
			k := datastore.NameKey(ItemKind, TestItem.ID, nil)
			var got Item
			err := client.Get(ctx, k, &got)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("retrieved: %s %v\n", got.ID, got.Title)
		})
	}
}
