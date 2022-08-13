.PHONY: run test import

run:
	@node serve.js

test:
	@echo
	@echo ===============================================
	@echo Backend tests
	@echo ===============================================
	@echo
	@cargo test

	@echo
	@echo ===============================================
	@echo Frontend tests
	@echo ===============================================
	@echo
	@npm test

	@echo
	@echo ===============================================
	@echo Import tests
	@echo ===============================================
	@echo
	@go test ./import/...

import:
	@go run import/main.go
