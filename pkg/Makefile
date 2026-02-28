SHELL := /bin/bash
.PHONY: all
all:
	@rm -f rulemancer/lex.go
	@re2c --lang go -o rulemancer/lex.go rulemancer/lex.go.re -i --api simple

.PHONY: clean
clean:
	@rm -f rulemancer/lex.go