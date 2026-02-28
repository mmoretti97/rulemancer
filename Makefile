
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
	@make -C docs --no-print-directory clean
	@make -C pkg --no-print-directory clean
