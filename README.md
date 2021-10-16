# link-checker
Website internal and external link checker in Go.

## Usage

1. Rename `config-sample.yaml` to `config.yaml`
2. Install dependencies:`go mod download`
3. Run `go run cmd/link-checker.go PROFILE`. E.g. `go run cmd/link-checker.go life`
4. When finished, a html report file will be generated under `reports/`
5. See `config.yaml` for details

Copyright 2021 Oleg Fomin, [ofstudio@gmail.com](mailto:ofstudio@gmail.com)

MIT License
