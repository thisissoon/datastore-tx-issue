package datastoretest

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// Cleanup deletes all entries for a kind
func Cleanup(ctx context.Context, t testing.TB, c *datastore.Client, kind string) {
	defer c.Close()
	q := datastore.NewQuery(kind)
	it := c.Run(ctx, q)
	for {
		var k *datastore.Key
		var err error
		k, err = it.Next(nil)
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatalf("Error fetching next item: %v", err)
		}
		if k == nil {
			break
		}
		t.Log("deleting key", k)
		err = c.Delete(ctx, k)
		if err != nil {
			t.Fatalf("error deleting item: %v", err)
		}
	}
}
