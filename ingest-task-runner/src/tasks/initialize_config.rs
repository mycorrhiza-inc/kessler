use async_trait::async_trait;

use crate::{
    common::tasks::{ExecuteUserTask, display_error_as_json},
    types::s3_stuff::{DIGITALOCEAN_S3, OPENSCRAPERS_S3_BUCKET},
};

use schemars::JsonSchema;
use serde::Deserialize;

#[derive(Debug, Copy, Clone, Deserialize, JsonSchema)]
pub struct InitializeConfig {}
#[async_trait]
impl ExecuteUserTask for InitializeConfig {
    async fn execute_task(self: Box<Self>) -> Result<serde_json::Value, serde_json::Value> {
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
                "Effect": "Allow",
                "Principal": "*",
                "Action": "s3:GetObject",
                "Resource": [
                    format!("arn:aws:s3:::{}/*", bucket)
                ]
            }
        ]
    })
    .to_string();
    let res = digitalocean_client
        .put_bucket_policy()
        .bucket(bucket)
        .policy(policy)
        .send()
        .await;
    match res {
        Ok(_) => Ok(()),
        Err(err) => Err(display_error_as_json(&err)),
    }
}
