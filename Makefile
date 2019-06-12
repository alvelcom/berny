.PHONY: build
build:
	@mkdir -p build
	go build -o build/berny cmd/berny/*.go
	go build -o build/bernyd cmd/bernyd/*.go

.PHONY: test unit-test e2e-test
test: unit-test e2e-test
unit-test:
	go test ./pkg/... ./cmd/...
e2e-test:
	ginkgo test/e2e
