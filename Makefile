.PHONY: run test

run:
	@node serve.js

test:
	@cargo test
	@npm test
