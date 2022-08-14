.PHONY: run test import test-import

run:
	@node serve.js

test: test-import
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

test-import:
	@echo
	@echo ===============================================
	@echo Import tests
	@echo ===============================================
	@echo
	@go test ./import/...

import:
	@go run import/main.go -output=data
