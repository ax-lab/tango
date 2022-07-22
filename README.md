# tango

Japanese dictionary and vocabulary manager.

## Development

This project requires a functional Node.js and Rust environment.

### Font-end

- Requires a Node.js environment.
- Use `npm install` to set up for the first time.
- Once set up, run with `npm start`.
- Note that the back-end must be run concurrently.

### Back-end

The server requires `cargo-watch` and `systemfd`. To install them use:

```sh
cargo install systemfd cargo-watch
```

Start the server with `make serve`.
