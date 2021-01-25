build:
	go build -o bin/main src/main.go src/restApi.go src/runningStats.go

run:
	go run src/main.go src/restApi.go src/runningStats.go

test:
	go test ./src/

clean:
	rm -rf src/main
