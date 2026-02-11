
all: rulemancer

.PHONY: rulemancer
rulemancer:
	@make -C pkg --no-print-directory all
	@go build
	@./rulemancer build

.PHONY: clean
clean:
	@rm -f ./rulemancer
	@rm -rf ./interface
