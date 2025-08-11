use std::convert::identity;

use aide::axum::ApiRouter;
use nypuc_ingest::NyPucIngestFull;
use transfer_files::TransferOpenscraperFilesIntoSupabase;

use crate::common::tasks::routing::declare_task_route;

pub mod initialize_config;
pub mod nypuc_ingest;
pub mod transfer_files;

pub fn add_user_task_routes(router: ApiRouter) -> ApiRouter {
    let router = declare_task_route::<NyPucIngestFull>(router);
    let router = declare_task_route::<TransferOpenscraperFilesIntoSupabase>(router);

    identity(router)
}
