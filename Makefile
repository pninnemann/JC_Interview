build:
    go build src/main.go src/restApi.go src/runningStats.go

test:
    go test

clean:
    rm -rf src/main