use std::{path::Path, str::FromStr, sync::LazyLock};

use anyhow::anyhow;
use async_trait::async_trait;
use aws_sdk_s3::{Client as S3Client, primitives::ByteStream};

use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;
use tracing::{debug, error, info, instrument};

use crate::{
    common::{
        hash::Blake2bHash,
        s3_generic::{S3Credentials, S3EnvNames, make_s3_lazylock},
        tasks::ExecuteUserTask,
    },
    types::openscrapers::{JurisdictionInfo, RawAttachment},
};

#[derive(Clone, Default, Deserialize, JsonSchema)]
pub struct TransferOpenscraperFilesIntoSupabase {
    only_transfer: Option<Vec<JurisdictionInfo>>,
    digitalocean_source_bucket: String,
    supabase_destination_bucket: String,
}

#[async_trait]
impl ExecuteUserTask for TransferOpenscraperFilesIntoSupabase {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = transfer_s3_files_to_supabase(
            self.only_transfer.as_deref(),
            &self.digitalocean_source_bucket,
            &self.supabase_destination_bucket,
        )
        .await;
        match res {
            Ok(()) => Ok("Task Completed Successfully".into()),
            Err(err) => Err(err.to_string().into()),
        }
    }
    fn get_task_label(&self) -> &'static str {
        "transfer_s3_files"
    }
    fn get_task_label_static() -> &'static str
    where
        Self: Sized,
    {
        "transfer_s3_files"
    }
}

struct SupS3 {}
impl S3EnvNames for SupS3 {
    const REGION_ENV: &str = "SUPABASE_S3_REGION";
    const ENDPOINT_ENV: &str = "SUPABASE_S3_ENDPOINT";
    const ACCESS_ENV: &str = "SUPABASE_S3_ACCESS_KEY";
    const SECRET_ENV: &str = "SUPABASE_S3_SECRET_KEY";
}
static SUPABASE_S3: LazyLock<S3Credentials> = make_s3_lazylock::<SupS3>();

struct OceanS3 {}
impl S3EnvNames for OceanS3 {
    const REGION_ENV: &str = "DIGITALOCEAN_S3_REGION";
    const ENDPOINT_ENV: &str = "DIGITALOCEAN_S3_ENDPOINT";
    const ACCESS_ENV: &str = "DIGITALOCEAN_S3_ACCESS_KEY";
    const SECRET_ENV: &str = "DIGITALOCEAN_S3_SECRET_KEY";
}

static DIGITALOCEAN_S3: LazyLock<S3Credentials> = make_s3_lazylock::<OceanS3>();

async fn transfer_s3_files_to_supabase(
    only_transfer: Option<&[JurisdictionInfo]>,
    source_bucket: &str,
    target_bucket: &str,
) -> anyhow::Result<()> {
    use futures::stream::{self, StreamExt};
    const RAW_FILE_PREFIX: &str = "files/raw/";
    let s3_supabase = SUPABASE_S3.make_s3_client().await;
    let s3_ocean = DIGITALOCEAN_S3.make_s3_client().await;

    let transfer_hashset = only_transfer.map(|jurisdictions| {
        jurisdictions
            .iter()
            .cloned()
            .collect::<std::collections::HashSet<JurisdictionInfo>>()
    });

    let all_hashes: Vec<Blake2bHash> = s3_ocean
        .list_objects_v2()
        .bucket(source_bucket)
        .prefix(RAW_FILE_PREFIX)
        .into_paginator()
        .send()
        .collect::<Vec<_>>()
        .await
        .into_iter()
        .filter_map(|result| match result {
            Ok(output) => Some(output),
            Err(err) => {
                error!(error = %err, "Error listing objects");
                None
            }
        })
        .flat_map(|output| output.contents.unwrap_or_default())
        .filter_map(|object| {
            object.key.and_then(|key| {
                Path::new(&key)
                    .file_name()
                    .and_then(|s| s.to_str())
                    .and_then(|s| Blake2bHash::from_str(s).ok())
            })
        })
        .collect();
    info!(hashes_count = %all_hashes.len(),"Collected hashes from s3 to try and attempt transfer");

    let should_be_transfered_over = |value: &JurisdictionInfo| -> bool {
        match &transfer_hashset {
            // In the case where no value is set, transfer over all files.
            None => true,
            Some(set) => set.contains(value),
        }
    };

    let process_hash = async |hash: Blake2bHash| {
        info!(%hash, "Processing file");

        let metadata_key = get_raw_attach_obj_key(hash);
        let metadata_bytes = match download_s3_bytes(&s3_ocean, source_bucket, &metadata_key).await
        {
            Ok(b) => b,
            Err(e) => {
                error!(%hash, error = %e, "Failed to download metadata, skipping");
                return;
            }
        };

        let raw_attachment: Option<RawAttachment> = match serde_json::from_slice(&metadata_bytes) {
            Ok(att) => Some(att),
            Err(e) => {
                error!(%hash, error = %e, "Failed to deserialize metadata, transferring anyway");
                None
            }
        };

        let should_transfer = match &raw_attachment {
            Some(att) => should_be_transfered_over(&att.jurisdiction_info),
            None => {
                debug!(%hash, "Transferring file despite metadata deserialization failure");
                true
            }
        };

        if should_transfer {
            info!(%hash, "Transfering file and metadata");
            let raw_file_key = get_raw_attach_file_key(hash);
            let raw_file_bytes =
                match download_s3_bytes(&s3_ocean, source_bucket, &raw_file_key).await {
                    Ok(b) => b,
                    Err(e) => {
                        error!(%hash, error = %e, "Failed to download raw file, skipping");
                        return;
                    }
                };

            if let Err(e) =
                upload_s3_bytes(&s3_supabase, target_bucket, &raw_file_key, raw_file_bytes).await
            {
                error!(%hash, error = %e, "Failed to upload raw file");
            }
            if let Err(e) =
                upload_s3_bytes(&s3_supabase, target_bucket, &metadata_key, metadata_bytes).await
            {
                error!(%hash, error = %e, "Failed to upload metadata");
            }
            info!(%hash, "Successfully transfered file and metadata");
        } else {
            debug!(%hash, "Skipping file, jurisdiction does not match");
        }
    };

    stream::iter(all_hashes.into_iter())
        .for_each_concurrent(10, process_hash)
        .await;

    Ok(())
}
// Core function to download bytes from S3
#[instrument(skip(s3_client))]
pub async fn download_s3_bytes(
    s3_client: &S3Client,
    bucket: &str,
    key: &str,
) -> anyhow::Result<Vec<u8>> {
    debug!(%bucket, %key,"Downloading S3 object");
    let output = s3_client
        .get_object()
        .bucket(bucket)
        .key(key)
        .send()
        .await
        .map_err(|e| {
            error!(error = %e, %bucket, %key,"Failed to download S3 object");
            e
        })?;

    let bytes = output
        .body
        .collect()
        .await
        .map(|data| data.into_bytes().to_vec())
        .map_err(|e| {
            error!(error = %e,%bucket, %key, "Failed to read response body");
            e
        })?;

    debug!(
        %bucket,
        %key,
        bytes_len = %bytes.len(),
        "Successfully downloaded file from s3"
    );
    Ok(bytes)
}

// Core function to upload bytes to S3
#[instrument(skip(s3_client, bytes))]
pub async fn upload_s3_bytes(
    s3_client: &S3Client,
    bucket: &str,
    key: &str,
    bytes: Vec<u8>,
) -> anyhow::Result<()> {
    debug!(len=%bytes.len(), %bucket, %key,"Uploading bytes to S3 object");
    s3_client
        .put_object()
        .bucket(bucket)
        .key(key)
        .body(ByteStream::from(bytes))
        .send()
        .await
        .map_err(|err| {
            error!(%err,%bucket, %key,"Failed to upload S3 object");
            anyhow!(err)
        })?;
    debug!( %bucket, %key,"Successfully uploaded s3 object");
    Ok(())
}

pub fn get_raw_attach_obj_key(hash: Blake2bHash) -> String {
    let key = format!("files/metadata/{hash}.json");
    debug!(%hash, "Generated raw attachment object key: {}", key);
    key
}

pub fn get_raw_attach_file_key(hash: Blake2bHash) -> String {
    let key = format!("files/raw/{hash}");
    debug!(%hash, "Generated raw attachment file key: {}", key);
    key
}
