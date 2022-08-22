BIN?=crudeci
SRC?=job.go main.go pipeline.go step.go


.PHONY: build
build:
	go build -o $(BIN) $(SRC)

run: build
	./$(BIN)

test:
	go test -v
