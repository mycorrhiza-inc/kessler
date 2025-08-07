#![allow(dead_code)]

use std::{
    convert::Infallible,
    net::{Ipv4Addr, SocketAddr},
};

use aide::axum::ApiRouter;
use common::{
    api_documentation::generate_api_docs_and_serve,
    otel_tracing::initialize_tracing_and_wrap_router,
    task_workers::{define_generic_task_routes, spawn_worker_loop},
};

mod common;
mod tasks;
mod types;
#[tokio::main]
async fn main() -> anyhow::Result<Infallible> {
    // initialise our subscriber
    let app_maker = || define_generic_task_routes(ApiRouter::new());
    // Add HTTP tracing layer
    // include trace context as header into the response

    let app = initialize_tracing_and_wrap_router(app_maker)?;
    // Spawn background worker to process PDF tasks
    // This worker runs indefinitely
    spawn_worker_loop();

    // bind and serve
    let addr = SocketAddr::new(Ipv4Addr::UNSPECIFIED.into(), 8123);
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    let Err(serve_err) = generate_api_docs_and_serve(listener, app, "A PDF processing API").await;
    Err(serve_err.into())
}
