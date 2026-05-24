clean:
	-rm -r build

get-dependencies:
	go get -v -t ./...

fmt:
	$(info Reformatting all source files...)
	go fmt ./...

build: clean fmt get-dependencies
	go build -o build/organize go.roman.zone/photo-tools/cmd/organize
	go build -o build/xmp-cleanup go.roman.zone/photo-tools/cmd/xmp-cleanup
	go build -o build/xmp-rating-clean go.roman.zone/photo-tools/cmd/xmp-rating-clean
