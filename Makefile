build:
	go build -o asciizer .

run:
	go run .

test:
	go test -v

deploy: build install-man install-completion
	cp asciizer ~/.local/bin/

install-man:
	install -d /usr/local/share/man/man1
	install -m 644 asciizer.1 /usr/local/share/man/man1/asciizer.1

install-completion:
	install -d ~/.oh-my-zsh/custom/completions
	install -m 644 _asciizer ~/.oh-my-zsh/custom/completions/_asciizer
