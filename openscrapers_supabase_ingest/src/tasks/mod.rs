use aide::axum::{ApiRouter, routing::post};
use nypuc_ingest::add_nypuc_all_task;

pub mod nypuc_ingest;

pub fn add_user_task_routes(router: ApiRouter) -> ApiRouter {
    router.api_route("tasks/new/ingest_all_nypuc", post(add_nypuc_all_task))
}
