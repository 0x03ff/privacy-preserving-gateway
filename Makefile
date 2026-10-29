.PHONY: all wasm frontend build clean run

# Default target: builds everything end-to-end
all: wasm frontend build

# 1. Compile FAME to WebAssembly
wasm:
	@echo "--- [1/3] Building FAME WebAssembly ---"
	cd scripts/fame && ./script.sh

# 2. Build SvelteKit Frontend
frontend:
	@echo "--- [2/3] Installing Node dependencies & building SvelteKit ---"
	npm install
	npm run build

# 3. Compile Go Backend Binary
build:
	@echo "--- [3/3] Compiling Go Backend Binary ---"
	go build -o bin/gateway ./cmd/main.go

# Clean up build artifacts and local storage
clean:
	rm -rf bin/
	rm -rf build/
	rm -rf .svelte-kit/
	rm -rf static/wasm/fame.wasm
	rm -rf static/wasm/wasm_exec.js

# Run the backend server (builds everything first)
run: all
	./bin/gateway