TEST?=$$(go list ./... | grep -v 'vendor')

up:
	docker-compose up

build:
	mkdir -p bin
	go build -o bin/guestlist main.go

test:
	docker-compose up -d db
	sleep 5
	go clean -testcache
	go test -v -i $(TEST) || exit 1
		echo $(TEST) | \
		xargs -t -n1 go test $(TESTARGS) -timeout=600s
	docker-compose stop db
