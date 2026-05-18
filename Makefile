.PHONY: test fmt run-chapter list

CH ?= 01

test:
	go test ./...

fmt:
	go fmt ./...

list:
	go run ./cmd/chapter

run-chapter:
	go run ./cmd/chapter -chapter $(CH)
