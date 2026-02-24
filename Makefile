GO = go
BINARY = ws
INSTALL_DIR = $(HOME)/.local/bin

.PHONY: build install clean

build:
	$(GO) build -o $(BINARY) .

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)

clean:
	rm -f $(BINARY)
