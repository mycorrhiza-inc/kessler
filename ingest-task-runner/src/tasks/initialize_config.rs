use async_trait::async_trait;
use mycorrhiza_common::tasks::{ExecuteUserTask, display_error_as_json};

use crate::types::s3_stuff::{DIGITALOCEAN_S3, OPENSCRAPERS_S3_BUCKET};

use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

#[derive(Debug, Copy, Clone, Serialize, Deserialize, JsonSchema, Default)]
pub struct InitializeConfig {}
#[async_trait]
impl ExecuteUserTask for InitializeConfig {
    async fn execute_task(self: Box<Self>) -> Result<serde_json::Value, serde_json::Value> {
        test_s3_client_permissions().await;
        make_openscrapers_public().await?;
        Ok("Successfully Configured Everything".into())
    }
    fn get_task_label_static() -> &'static str
    where
        Self: Sized,
    {
        "initialize_config"
    }
    fn get_task_label(&self) -> &'static str {
        "initialize_config"
    }
}

async fn make_openscrapers_public() -> Result<(), serde_json::Value> {
    let digitalocean_client = DIGITALOCEAN_S3.make_s3_client().await;
    let bucket = &**OPENSCRAPERS_S3_BUCKET;
    let policy = serde_json::json!({
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "PublicReadGetObject",
                "Effect": "Allow",
                // "Principal": {"AWS": "*"}, // Try this format instead of just "*"
                "Principal": "*", // Use just "*" for DigitalOcean Spaces, not {"AWS": "*"}
                "Action": "s3:GetObject",
                // "Resource": format!("arn:aws:s3:::{}/*", bucket),
                "Resource": format!("arn:digitalocean:spaces:*:*:{}/", bucket),
            }
        ]
    });
    let policy_string = policy.to_string();
    tracing::info!(bucket, policy_string, "Saving policy to s3");

    let res = digitalocean_client
        .put_bucket_policy()
        .bucket(bucket)
        .policy(policy_string)
        .send()
        .await;
    match res {
        Ok(_) => Ok(()),
        Err(err) => Err(display_error_as_json(&err)),
    }
}

use tracing::info;

/// Checks that the S3 key‑pair we are using has enough permissions to
/// perform the most common bucket operations.
pub async fn test_s3_client_permissions() {
    // Initialise the client once – it is reused for all requests.
    let digitalocean_client = DIGITALOCEAN_S3.make_s3_client().await;
    // `OPENSCRAPERS_S3_BUCKET` is a `Lazy<String>`/`Arc<String>` in the original code,
    // `&**` dereferences it to a plain `&str`.
    let bucket = &**OPENSCRAPERS_S3_BUCKET;

    // -------------------------------------------------------------------------
    // 1️⃣  List buckets
    // -------------------------------------------------------------------------
    info!(action = "list_buckets", bucket = %bucket, "starting list_buckets test");
    match digitalocean_client.list_buckets().send().await {
        Ok(resp) => {
            // `resp.buckets()` returns an `Option<&[Bucket]>`
            let bucket_cnt = resp.buckets().len();
            info!(
                action = "list_buckets",
                status = "success",
                buckets = bucket_cnt,
                "list_buckets succeeded"
            );
        }
        Err(err) => {
            info!(
                action = "list_buckets",
                status = "failed",
                error = %err,
                "list_buckets failed"
            );
        }
    }

    // -------------------------------------------------------------------------
    // 2️⃣  Get bucket location
    // -------------------------------------------------------------------------
    info!(action = "get_bucket_location", bucket = %bucket, "starting get_bucket_location test");
    match digitalocean_client
        .get_bucket_location()
        .bucket(bucket)
        .send()
        .await
    {
        Ok(_) => {
            info!(
                action = "get_bucket_location",
                status = "success",
                bucket = %bucket,
                "get_bucket_location succeeded"
            );
        }
        Err(err) => {
            info!(
                action = "get_bucket_location",
                status = "failed",
                bucket = %bucket,
                error = %err,
                "get_bucket_location failed"
            );
        }
    }

    // -------------------------------------------------------------------------
    // 3️⃣  Get bucket ACL
    // -------------------------------------------------------------------------
    info!(action = "get_bucket_acl", bucket = %bucket, "starting get_bucket_acl test");
    match digitalocean_client
        .get_bucket_acl()
        .bucket(bucket)
        .send()
        .await
    {
        Ok(_) => {
            info!(
                action = "get_bucket_acl",
                status = "success",
                bucket = %bucket,
                "get_bucket_acl succeeded"
            );
        }
        Err(err) => {
            info!(
                action = "get_bucket_acl",
                status = "failed",
                bucket = %bucket,
                error = %err,
                "get_bucket_acl failed"
            );
        }
    }

    // -------------------------------------------------------------------------
    // 4️⃣  Get bucket policy
    // -------------------------------------------------------------------------
    info!(action = "get_bucket_policy", bucket = %bucket, "starting get_bucket_policy test");
    match digitalocean_client
        .get_bucket_policy()
        .bucket(bucket)
        .send()
        .await
    {
        Ok(_) => {
            info!(
                action = "get_bucket_policy",
                status = "success",
                bucket = %bucket,
                "bucket policy exists"
            );
        }
        Err(err) => {
            info!(
                action = "get_bucket_policy",
                status = "failed",
                bucket = %bucket,
                error = %err,
                "get_bucket_policy failed"
            );
        }
    }

    // -------------------------------------------------------------------------
    // 5️⃣  Put a simple object
    // -------------------------------------------------------------------------
    const TEST_KEY: &str = "test-permissions.txt";
    info!(
        action = "put_object",
        bucket = %bucket,
        key = TEST_KEY,
        "starting put_object test"
    );
    match digitalocean_client
        .put_object()
        .bucket(bucket)
        .key(TEST_KEY)
        .body("test".to_string().into_bytes().into())
        .send()
        .await
    {
        Ok(_) => {
            info!(
                action = "put_object",
                status = "success",
                bucket = %bucket,
                key = TEST_KEY,
                "put_object succeeded"
            );
        }
        Err(err) => {
            info!(
                action = "put_object",
                status = "failed",
                bucket = %bucket,
                key = TEST_KEY,
                error = %err,
                "put_object failed"
            );
        }
    }
}
