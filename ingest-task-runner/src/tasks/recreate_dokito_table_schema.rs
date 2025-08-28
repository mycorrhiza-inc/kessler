use async_trait::async_trait;
use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;
use sqlx::{Executor, PgConnection, PgPool, migrate::Migrator, postgres::PgPoolOptions};
use tracing::info;

use mycorrhiza_common::tasks::ExecuteUserTask;

use super::nypuc_ingest::DEFAULT_POSTGRES_CONNECTION_URL;

#[derive(Clone, Copy, Default, Deserialize, JsonSchema)]
pub struct RecreateDokitoTableSchema {}

#[async_trait]
impl ExecuteUserTask for RecreateDokitoTableSchema {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = recreate_schema().await;
        match res {
            Ok(()) => {
                info!("Recreated schema.");
                Ok("Task Completed Successfully".into())
            }
            Err(err) => {
                tracing::error!(error= % err, error_debug= ?err,"Encountered error in recreate_schema");
                Err(err.to_string().into())
            }
        }
    }
    fn get_task_label(&self) -> &'static str {
        "recreate_dokito_table_schema"
    }
    fn get_task_label_static() -> &'static str
    where
        Self: Sized,
    {
        "recreate_dokito_table_schema"
    }
}

pub async fn recreate_schema() -> anyhow::Result<()> {
    info!("Got request to recreate schema");
    let db_url = &**DEFAULT_POSTGRES_CONNECTION_URL;
    let pool = PgPoolOptions::new().connect(db_url).await?;
    info!("Created pg pool");

    let mut migrator = sqlx::migrate!("./complete_migrations");

    let num_migrations = migrator.iter().count();
    info!(%num_migrations,"Created sqlx migrator");

    info!("Dropping existing tables");
    drop_existing_schema(&pool, &mut migrator).await?;

    info!("Creating tables");
    // create_schema(&pool).await?;
    create_schema(&pool, &mut migrator).await?;

    info!("Successfully recreated schema");

    Ok(())
}

pub async fn drop_existing_schema(pool: &PgPool, _migrator: &mut Migrator) -> anyhow::Result<()> {
    // migrator.set_ignore_missing(true).undo(pool, 0).await?;
    pool.execute(include_str!(
        "../../complete_migrations/001_dokito_complete.down.sql"
    ))
    .await?;
    Ok(())
}

pub async fn create_schema(pool: &PgPool, _migrator: &mut Migrator) -> anyhow::Result<()> {
    // migrator.set_ignore_missing(true).run(pool).await?;

    pool.execute(include_str!(
        "../../complete_migrations/001_dokito_complete.up.sql"
    ))
    .await?;
    Ok(())
}
