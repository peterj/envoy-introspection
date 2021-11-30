wasm:
	tinygo build -o envoy-introspection.wasm -scheduler=none -target=wasi main.go