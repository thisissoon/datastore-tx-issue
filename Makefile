start-ds:
	gcloud beta emulators datastore start --consistency=1

init-ds:
	gcloud beta emulators datastore env-init

test:
	go test -v
