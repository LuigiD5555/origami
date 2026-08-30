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
	go build -trimpath -o bin/origami-hybrid-carrier ./cmd/origami-hybrid-carrier
	go build -trimpath -o bin/origami-hybrid-tool ./cmd/origami-hybrid-tool
	go build -trimpath -o bin/origami-memory ./cmd/origami-memory
	go build -trimpath -o bin/origami-memory-scale ./cmd/origami-memory-scale
	go build -trimpath -o bin/origami-fixed-carrier ./cmd/origami-fixed-carrier
	go build -trimpath -o bin/origami-reduce ./cmd/origami-reduce

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
