use std::{env, sync::LazyLock};

use mycorrhiza_common::s3_generic::{S3Credentials, S3EnvNames, make_s3_lazylock};

struct SupS3 {}
impl S3EnvNames for SupS3 {
    const REGION_ENV: &str = "SUPABASE_S3_REGION";
    const ENDPOINT_ENV: &str = "SUPABASE_S3_ENDPOINT";
    const ACCESS_ENV: &str = "SUPABASE_S3_ACCESS_KEY";
    const SECRET_ENV: &str = "SUPABASE_S3_SECRET_KEY";
}

pub static SUPABASE_S3: LazyLock<S3Credentials> = make_s3_lazylock::<SupS3>();
pub static OPENSCRAPERS_S3_BUCKET: LazyLock<String> =
    LazyLock::new(|| env::var("OPENSCRAPERS_S3_BUCKET").unwrap_or("openscrapers".to_string()));

struct OceanS3 {}
impl S3EnvNames for OceanS3 {
    const REGION_ENV: &str = "DIGITALOCEAN_S3_REGION";
    const ENDPOINT_ENV: &str = "DIGITALOCEAN_S3_ENDPOINT";
    const ACCESS_ENV: &str = "DIGITALOCEAN_S3_ACCESS_KEY";
    const SECRET_ENV: &str = "DIGITALOCEAN_S3_SECRET_KEY";
}

pub static DIGITALOCEAN_S3: LazyLock<S3Credentials> = make_s3_lazylock::<OceanS3>();
