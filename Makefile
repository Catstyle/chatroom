.PHONY: _default clean build test run

_default: build

clean:
	rm -r ./bin

build:
	GOARCH=amd64 go build -ldflags "-s -w" -a -o bin/server .

test:
	go test ./...

run:
	./bin/server

.PHONY: list
list:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
