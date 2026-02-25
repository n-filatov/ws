GO = go
BINARY = ws
INSTALL_DIR = $(HOME)/.local/bin

.PHONY: build install clean snapshot release

build:
	$(GO) build -o $(BINARY) .

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)

clean:
	rm -f $(BINARY)

# Build a local snapshot (no release, no push) — useful for testing goreleaser config
snapshot:
	goreleaser release --snapshot --clean

# Tag and push to trigger the release CI. Usage: make release VERSION=v0.1.0
release:
	@test -n "$(VERSION)" || (echo "Usage: make release VERSION=v0.1.0" && exit 1)
	git tag $(VERSION)
	git push origin $(VERSION)
