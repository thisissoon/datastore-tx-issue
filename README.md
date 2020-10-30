# Datastore Concurrent Transaction Error

This repo demonstrates an issue we have been seeing with Google Datastore.

When using concurrent transactions we get this error `datastore: concurrent transaction`. According to the [docs](https://cloud.google.com/datastore/docs/concepts/transactions#transaction_locks):

> When two concurrent read-write transactions read or write the same data, the lock held by one transaction can delay the other transaction.

This suggests that datastore will lock a transaction to prevent concurrent read/writes by a separate transaction running at the same time, but this is not what we are seeing here.

To run this example you will need the datastore emulator. Instructions for installing this can be found [here](https://cloud.google.com/datastore/docs/tools/datastore-emulator).

This project uses Go v1.15.3 but should work with older versions

## Testing

To reproduce the error this example uses testing. The test file is main_test.go and the test is `TestStore_PutTx` This will concurrently call the method `store.PutTx` x number of times where x is the number of putItems in the table test. There are two test cases:
1. `no lock` will fail most of the time as it relies on datastore locking concurrent transactions which it doesn't appear to be doing.
2. `with lock` uses our own locking solution with sync.Mutex and will pass all of the time.

To run the tests
1. run `make start-ds`. This will start the emulator with the maximum amount of consistency
2. In a separate terminal run `make init-ds`. This will set environment variables for communicating with the emulator
3. Run `make test`
