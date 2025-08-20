use std::convert::identity;

use aide::axum::ApiRouter;
use nypuc_ingest::NyPucIngestPurgePrevious;
use transfer_files::TransferOpenscraperFilesIntoSupabase;

use crate::{
    common::tasks::routing::declare_task_route,
    tasks::{initialize_config::InitializeConfig, nypuc_ingest::NyPucIngestGetMissingDockets},
};

pub mod initialize_config;
pub mod nypuc_ingest;
pub mod transfer_files;

pub fn add_user_task_routes(router: ApiRouter) -> ApiRouter {
    let router = declare_task_route::<NyPucIngestPurgePrevious>(router);
    let router = declare_task_route::<NyPucIngestGetMissingDockets>(router);
    let router = declare_task_route::<TransferOpenscraperFilesIntoSupabase>(router);
    let router = declare_task_route::<InitializeConfig>(router);

    identity(router)
}
