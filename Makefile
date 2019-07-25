bin/migrate-actions: $(go_src)
	script/build

unit:
	go test -v ./...

test: unit bin/migrate-actions
	script/fixtures-test

