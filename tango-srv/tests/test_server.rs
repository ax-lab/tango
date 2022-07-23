mod common;
use std::error::Error;

use assert_json_diff::*;
use common::*;

use serde_json::json;
use testable::*;

#[test]
fn server_outputs_bound_port_on_start() {
	let (mut cleanup, get_port) = spawn_server();
	defer!(cleanup());

	let port = panic_after(Duration::from_millis(50), get_port);
	let port = port.expect("server did not output port number");
	assert!(port > 0);
}

#[test]
fn server_accepts_connections() {
	let response = test_request(Method::GET, "api/ping");
	assert_eq!(response.status(), 200);
	assert_eq!(response.text().unwrap(), "pong");
}

#[test]
fn server_accepts_graphql() -> Result<(), Box<dyn Error>> {
	let response = test_request_with_config(Method::POST, "api/query", |request| {
		let headers = request.headers_mut();
		headers.append(CONTENT_TYPE, "application/graphql".parse().unwrap());
		request.body_mut().replace(Body::from(r"query { appName }"));
	});
	assert_eq!(response.status(), 200);
	assert!(response
		.headers()
		.get(CONTENT_TYPE)
		.unwrap()
		.to_str()?
		.contains("application/json"));
	let response = response.text()?;
	let response = serde_json::from_str::<serde_json::Value>(&response)?;
	let expected = json!({ "data": { "appName": "Tango" } });
	assert_json_include!(actual: response, expected: expected);

	Ok(())
}

#[test]
fn server_accepts_graphiql() -> Result<(), Box<dyn Error>> {
	let response = test_request(Method::GET, "api/graphql");
	assert_eq!(response.status(), 200);
	assert!(response
		.headers()
		.get(CONTENT_TYPE)
		.unwrap()
		.to_str()?
		.contains("text/html"));
	let response = response.text()?;
	assert!(response.contains("GraphiQL"));

	Ok(())
}
