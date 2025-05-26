#data:
#	cp files/* static/
compile: #data
	GOOS=js GOARCH=wasm go build -o static/main.wasm
serve: compile
	go run http/main.go