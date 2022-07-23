pub use std::time::Duration;

use reqwest::{blocking::*, Url};
use testable::*;

pub use reqwest::blocking::Body;
pub use reqwest::header::CONTENT_TYPE;
pub use reqwest::Method;

pub fn test_request(method: Method, path: &str) -> Response {
	test_request_with_config(method, path, |_| {})
}

pub fn test_request_with_config<F: FnOnce(&mut Request)>(
	method: Method,
	path: &str,
	config: F,
) -> Response {
	let (mut cleanup, get_port) = spawn_server();
	defer!(cleanup());

	let port =
		panic_after(Duration::from_millis(100), get_port).expect("server did not return a port");

	let url = Url::parse(&format!("http://127.0.0.1:{}/{}", port, path)).unwrap();
	let mut request = Request::new(method, url);
	config(&mut request);

	let client = ClientBuilder::new()
		.timeout(Duration::from_millis(200))
		.build()
		.unwrap();
	client.execute(request).unwrap()
}

pub fn spawn_server() -> (impl FnMut(), impl FnOnce() -> Option<u16>) {
	use std::io::*;
	use std::process::*;

	use regex::Regex;

	let mut cmd = tux::get_bin("tango-srv");
	let mut child = cmd
		.stdout(Stdio::null())
		.stderr(Stdio::piped())
		.spawn()
		.unwrap();
	let output = child.stderr.take().unwrap();
	let re = Regex::new(r"server-start-port=(\d+)").unwrap();
	let get_port = move || {
		let output = BufReader::new(output);
		for line in output.lines() {
			let line = line.expect("read output line failed");
			if let Some(captures) = re.captures(&line) {
				let port = captures.get(1).unwrap().as_str();
				let port = port.parse::<u16>().unwrap();
				return Some(port);
			}
		}
		None
	};

	let kill_child = move || {
		let _ = child.kill();
	};

	(kill_child, get_port)
}
