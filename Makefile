go       := go
out      := ./out
artifact := $(out)/bank-host

.PHONY: build
build:
	@mkdir -p $(out)
	$(go) build -o $(artifact)

.PHONY: run
run: build
	@echo "===== Running Bank =====\n"
	@$(artifact)

.PHONY: clean
clean:
	rm -rf $(out)