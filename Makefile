build:
	go build -ldflags="-s -w" -o link-checker ./cmd/link-checker.go

# To enable cross-compiling for CGO you need to have a local toolchain that can compile C code for that target.
# On macOS, you can install mingw with homebrew:
#   brew install mingw-w64
#
# https://stackoverflow.com/a/36916044
build_win64:
	CC="x86_64-w64-mingw32-gcc" CGO_ENABLED=1 \
	GOOS=windows GOARCH=amd64 \
	go build -ldflags="-s -w" -o link-checker.exe ./cmd/link-checker.go
