use std::sync::LazyLock;

use aws_config::{BehaviorVersion, Region};
use aws_sdk_s3::{Client, config::Credentials};
use tracing::info;

pub struct S3Credentials {
    cloud_region: String,
    endpoint: String,
    access_key: String,
    secret_key: String,
}

impl S3Credentials {
    pub async fn make_s3_client(&self) -> Client {
        info!("Creating S3 client");
        let creds = Credentials::new(
            &self.access_key,
            &self.secret_key,
            None, // no session token
            None, // no expiration
            "manual",
        );

        // Start from the env-loader so we still pick up other settings (timeouts, retry, etc)
        let cfg_loader = aws_config::defaults(BehaviorVersion::latest())
            .region(Region::new(self.cloud_region.clone()))
            .credentials_provider(creds)
            .endpoint_url(&self.endpoint);

        let sdk_config = cfg_loader.load().await;
        Client::new(&sdk_config)
    }
}

const DEFAULT_S3_REGION: &str = "sfo3";
const DEFAULT_S3_ENDPOINT: &str = "https://sfo3.digitaloceanspaces.com";

pub trait S3EnvNames {
    const REGION_ENV: &str;
    const ENDPOINT_ENV: &str;
    const ACCESS_ENV: &str;
    const SECRET_ENV: &str;
}
fn init_from_env_vars<T: S3EnvNames>() -> S3Credentials {
    let access_key = std::env::var(T::ACCESS_ENV).expect(&format!(
        "S3 access key env var was not set: {}",
        T::ACCESS_ENV
    ));
    let cloud_region = std::env::var(T::REGION_ENV).unwrap_or_else(|_| {
        println!("S3 region not set, using default: {DEFAULT_S3_REGION}");
        DEFAULT_S3_REGION.to_string()
    });
    let endpoint = std::env::var(T::ENDPOINT_ENV).unwrap_or_else(|_| {
        println!("S3 endpoint not set, using default: {DEFAULT_S3_ENDPOINT}");
        DEFAULT_S3_ENDPOINT.to_string()
    });
    let secret_key = std::env::var(T::SECRET_ENV).expect(&format!(
        "S3 secret key env var was not set: {}",
        T::SECRET_ENV
    ));

    S3Credentials {
        cloud_region,
        endpoint,
        access_key,
        secret_key,
    }
}
pub const fn s3_locked<T: S3EnvNames>() -> LazyLock<S3Credentials> {
    LazyLock::new(init_from_env_vars::<T>)
}
