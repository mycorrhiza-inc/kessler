
use aide::{
    axum::{ApiRouter, IntoApiResponse, routing::get_with},
    transform::TransformOperation,
};
use axum::{Json, extract::Path, response::IntoResponse};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use super::{TaskStatusDisplay, workers::TASK_STATUS_DATA};
pub const CHECK_TASK_URL_LEAF: &str = "/tasks/check_status";

#[derive(Clone, Copy, Serialize, Deserialize, JsonSchema)]
pub struct TaskIDNumber {
    task_id: u64,
}
pub async fn check_task_status(
    Path(TaskIDNumber { task_id }): Path<TaskIDNumber>,
) -> impl IntoApiResponse {
    let read_guard = (*TASK_STATUS_DATA).read().await;
    let status = read_guard.get(&task_id).cloned();
    drop(read_guard);
    match status {
        None => (axum::http::StatusCode::NOT_FOUND, format!("Could not find task with task_id: {task_id}")).into_response()
, // Return a 404 error with an error string "Could not find task with that id"
        Some(status) => {
            let display_status: TaskStatusDisplay = status.into();
            // Return a 200 code with the status that wants to be returned.
(axum::http::StatusCode::OK, Json(display_status)).into_response()
        }
    }
}

pub fn check_task_status_docs(op: TransformOperation) -> TransformOperation {
    op.description("Fetch attachment data from S3.")
        .response::<200, Json<TaskStatusDisplay>>()
        .response_with::<404, String, _>(|res| {
            res.description("Could not find task with that task_id")
        })
}

pub fn define_generic_task_routes(router: ApiRouter) -> ApiRouter {
    let check_path = CHECK_TASK_URL_LEAF.to_string() + "/{task_id}";
    router.api_route(
        &check_path,
        get_with(check_task_status, check_task_status_docs),
    )
}

#[derive(Clone, Copy, Default, Serialize, Deserialize, JsonSchema)]
pub struct PriorityExtractor {
    pub priority: i32,
}
