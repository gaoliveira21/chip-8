build_wasm:
	GOOS=js GOARCH=wasm go build -o ./web/public/chip8.wasm ./web/web.go

devserver:
	go run ./web/server/devserver.go

start: build_wasm devserver
