.PHONY:

run:
	go run ./main.go
test:
	go test --count=1 ./...
test.v:
	go test --count=1 -v ./...