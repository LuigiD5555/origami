.PHONY: test vet check build install uninstall prepare-microisa inspect-dimensions clean-runs

test:
	go test ./...

vet:
	go vet ./...

check: test vet

build:
	mkdir -p bin
	go build -trimpath -o bin/ohf-lab ./cmd/ohf-lab
	go build -trimpath -o bin/ohf-glyphcalc ./cmd/ohf-glyphcalc

install:
	./install.sh

uninstall:
	./uninstall.sh

prepare-microisa:
	./scripts/microisa_prepare_local.sh

inspect-dimensions:
	./scripts/dimensional_inspect.sh

clean-runs:
	find runs -mindepth 1 -maxdepth 1 ! -name '.gitkeep' -exec rm -rf {} +
