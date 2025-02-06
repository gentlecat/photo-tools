clean:
	-rm -r build

get-dependencies:
	go get -v -t ./...

fmt:
	$(info Reformatting all source files...)
	go fmt ./...

build: clean fmt get-dependencies
	go build -o build/organize go.roman.zone/photo-tools/cmd/organize
