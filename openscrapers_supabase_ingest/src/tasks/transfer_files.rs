use std::sync::LazyLock;

use anyhow::anyhow;
use async_trait::async_trait;
use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;

use crate::{
    common::{
        s3_generic::{S3Credentials, S3EnvNames, s3_locked},
        tasks::ExecuteUserTask,
    },
    types::openscrapers::JurisdictionInfo,
};

#[derive(Clone, Default, Deserialize, JsonSchema)]
pub struct TransferOpenscraperFilesIntoSupabase {
    only_transfer: Option<Vec<JurisdictionInfo>>,
}

#[async_trait]
impl ExecuteUserTask for TransferOpenscraperFilesIntoSupabase {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = Err(anyhow!("not implemented"));
        match res {
            Ok(()) => Ok("Task Completed Successfully".into()),
            Err(err) => Err(err.to_string().into()),
        }
    }
    fn get_task_label(&self) -> &'static str {
        "ingest_nypuc_all"
    }
    fn get_task_label_static() -> &'static str
    where
        Self: Sized,
    {
        "ingest_nypuc_all"
    }
}

struct SupS3 {}
impl S3EnvNames for SupS3 {
    const REGION_ENV: &str = "SUPABASE_S3_REGION";
    const ENDPOINT_ENV: &str = "SUPABASE_S3_ENDPOINT";
    const ACCESS_ENV: &str = "SUPABASE_S3_ACCESS_KEY";
    const SECRET_ENV: &str = "SUPABASE_S3_SECRET_KEY";
}
static SUPABASE_S3: LazyLock<S3Credentials> = s3_locked::<SupS3>();

struct OceanS3 {}
impl S3EnvNames for OceanS3 {
    const REGION_ENV: &str = "DIGITALOCEAN_S3_REGION";
    const ENDPOINT_ENV: &str = "DIGITALOCEAN_S3_ENDPOINT";
    const ACCESS_ENV: &str = "DIGITALOCEAN_S3_ACCESS_KEY";
    const SECRET_ENV: &str = "DIGITALOCEAN_S3_SECRET_KEY";
}

static DIGITALOCEAN_S3: LazyLock<S3Credentials> = s3_locked::<OceanS3>();

async fn transfer_s3_files_to_supabase(
    only_transfer: Option<&[JurisdictionInfo]>,
    source_bucket: &str,
    target_bucket: &str,
) -> anyhow::Result<()> {
    const TRANSFER_ALL_WITH_ROOT: &str = "/files/";
    let s3_supabase = SUPABASE_S3.make_s3_client().await;
    let s3_ocean = DIGITALOCEAN_S3.make_s3_client().await;

    Ok(())
    // Files flow from ocean -> supabase
}
