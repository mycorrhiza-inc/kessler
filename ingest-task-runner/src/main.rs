#![allow(dead_code)]

use std::{
    convert::{Infallible, identity},
    net::{Ipv4Addr, SocketAddr},
};

use aide::axum::ApiRouter;
use mycorrhiza_common::{
    api_documentation::generate_api_docs_and_serve,
    llm_deepinfra::DEEPINFRA_API_KEY,
    otel_tracing::initialize_tracing_and_wrap_router,
    tasks::{routing::define_generic_task_routes, workers::spawn_worker_loop},
};
use tasks::add_user_task_routes;

use crate::types::s3_stuff::DIGITALOCEAN_S3;

mod tasks;
mod types;
#[tokio::main]
async fn main() -> anyhow::Result<Infallible> {
    let _ = *DEEPINFRA_API_KEY;
    let _ = *DIGITALOCEAN_S3;
    // initialise our subscriber
    let app_maker = || {
        let router = ApiRouter::new();
        let router = define_generic_task_routes(router);
        let router = add_user_task_routes(router);
        identity(router)
    };
    // Add HTTP tracing layer
    // include trace context as header into the response

    let (app, _guard) = initialize_tracing_and_wrap_router(app_maker)?;
    // Spawn background worker to process PDF tasks
    // This worker runs indefinitely
    spawn_worker_loop();

    // bind and serve
    let addr = SocketAddr::new(Ipv4Addr::UNSPECIFIED.into(), 8123);
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    let Err(serve_err) = generate_api_docs_and_serve(listener, app, "A PDF processing API").await;
    Err(serve_err.into())
}
