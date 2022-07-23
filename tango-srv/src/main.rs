use std::convert::Infallible;

use env_logger::Builder;
use listenfd::ListenFd;
use log::*;
use warp::hyper;
use warp::Filter;

use juniper::*;

pub type Schema = RootNode<'static, Query, Mutation, EmptySubscription<Context>>;

pub struct Context;

impl juniper::Context for Context {}

pub struct Query;

#[graphql_object(Context = Context)]
impl Query {
	/// Returns the application name.
	pub fn app_name() -> &'static str {
		"Tango"
	}
}

pub struct Mutation;

#[graphql_object(Context = Context)]
impl Mutation {
	/// Computes the answer to life, the universe, and everything.
	fn answer() -> i32 {
		42
	}
}

fn schema() -> Schema {
	Schema::new(Query, Mutation, EmptySubscription::new())
}

#[tokio::main]
async fn main() {
	let mut builder = Builder::default();
	builder.filter_level(LevelFilter::Info);
	builder.parse_default_env();
	builder.init();

	info!("starting up server...");

	let api = warp::path!("api" / ..);

	let state = warp::any().map(|| Context);
	let graphql = juniper_warp::make_graphql_filter(schema(), state.boxed());
	let graphql = warp::path("query").and(graphql);
	let graphiql = warp::get()
		.and(warp::path("graphql"))
		.and(juniper_warp::graphiql_filter("/api/query", None));

	let ping = warp::path!("ping").map(|| "pong");
	let hello_name = warp::path!("hello" / String).map(|name| format!("Hello, {}!", name));
	let hello_stranger = warp::path!("hello").map(|| "Hello, stranger!");
	let hello = api.and(
		ping.or(graphql)
			.or(graphiql)
			.or(hello_name)
			.or(hello_stranger),
	);

	let routes = hello;

	let service = warp::service(routes);
	let make_service = hyper::service::make_service_fn(|_| {
		let service = service.clone();
		async move { Ok::<_, Infallible>(service) }
	});

	// Accept a socket provided using systemfd, failing back to a random port
	// otherwise.
	let mut listenfd = ListenFd::from_env();
	let server = if let Some(listener) = listenfd.take_tcp_listener(0).unwrap() {
		hyper::Server::from_tcp(listener).unwrap()
	} else {
		hyper::Server::bind(&([127, 0, 0, 1], 0).into())
	};

	let server = server.serve(make_service);

	let addr = server.local_addr();
	info!("listening at {}...", addr);
	info!("server-start-port={}", addr.port());

	let server = server.with_graceful_shutdown(async {
		tokio::signal::ctrl_c()
			.await
			.unwrap_or_else(|err| error!("failed to bind ctrl+c listener: {}", err));
	});

	server.await.unwrap();
	info!("server shut down");
}
