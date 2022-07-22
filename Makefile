.PHONY: run serve

WATCH= -w tango-srv

run:
	cargo run

serve:
	systemfd --no-pid -s http::29801 -- cargo watch $(WATCH) -- cargo run -q
