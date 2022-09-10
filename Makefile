.PHONY: run test import test-import json

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
	@go test -count=1 ./import/...

import:
	@go run import/main.go -output=data

json:
	@go run export-json/main.go -output=data/json -import=data
