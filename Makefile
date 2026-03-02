VERSION ?= dev
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o asciizer .

run:
	go run .

test:
	go test -v

release: clean test
	@mkdir -p dist
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/asciizer . && \
		tar -czf dist/asciizer-$(VERSION)-linux-amd64.tar.gz -C dist asciizer && rm dist/asciizer
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/asciizer . && \
		tar -czf dist/asciizer-$(VERSION)-linux-arm64.tar.gz -C dist asciizer && rm dist/asciizer
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/asciizer . && \
		tar -czf dist/asciizer-$(VERSION)-darwin-amd64.tar.gz -C dist asciizer && rm dist/asciizer
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/asciizer . && \
		tar -czf dist/asciizer-$(VERSION)-darwin-arm64.tar.gz -C dist asciizer && rm dist/asciizer
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/asciizer.exe . && \
		cd dist && zip asciizer-$(VERSION)-windows-amd64.zip asciizer.exe && rm asciizer.exe

clean:
	rm -rf dist/
	rm -f asciizer

deploy: build install-man install-completion
	cp asciizer ~/.local/bin/

install-man:
	install -d /usr/local/share/man/man1
	install -m 644 asciizer.1 /usr/local/share/man/man1/asciizer.1

install-completion:
	install -d ~/.oh-my-zsh/custom/completions
	install -m 644 _asciizer ~/.oh-my-zsh/custom/completions/_asciizer

.PHONY: build run test release clean deploy install-man install-completion
