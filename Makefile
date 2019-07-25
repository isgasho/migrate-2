bin/migrate-actions: $(go_src)
	script/build

test: bin/migrate-actions
	go test -v ./... && script/fixtures-test

