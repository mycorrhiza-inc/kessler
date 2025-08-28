use async_trait::async_trait;
use futures::{StreamExt, stream};
use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;
use sqlx::postgres::PgPoolOptions;
use sqlx::types::Uuid;
use tracing::info;

use crate::tasks::nypuc_ingest::DEFAULT_POSTGRES_CONNECTION_URL;
use mycorrhiza_common::tasks::ExecuteUserTask;

#[derive(Clone, Copy, Default, Deserialize, JsonSchema)]
pub struct SplitCompatifiedTypeIntoSubtype {}

#[async_trait]
impl ExecuteUserTask for SplitCompatifiedTypeIntoSubtype {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = process_docket_types().await;
        match res {
            Ok(()) => {
                info!("Temp script task completed successfully");
                Ok("Temp script task completed successfully".into())
            }
            Err(err) => {
                tracing::error!(error = %err, "Encountered error in temp script task");
                Err(err.to_string().into())
            }
        }
    }

    fn get_task_label(&self) -> &'static str {
        "split_compatified_type_into_subtype"
    }

    fn get_task_label_static() -> &'static str
    where
        Self: Sized,
    {
        "split_compatified_type_into_subtype"
    }
}

async fn process_docket_types() -> anyhow::Result<()> {
    info!("Starting temp script task to process docket types");

    let db_url = &**DEFAULT_POSTGRES_CONNECTION_URL;
    let pool = PgPoolOptions::new()
        .max_connections(10)
        .connect(db_url)
        .await?;
    info!("Created pg pool");

    // Get all dockets where docket_type is not null and not empty,
    // and docket_subtype is null or empty
    let dockets = sqlx::query!(
        r#"
        SELECT uuid, docket_type 
        FROM dockets 
        WHERE docket_type IS NOT NULL 
        AND docket_type != '' 
        AND (docket_subtype IS NULL OR docket_subtype = '')
        "#
    )
    .fetch_all(&pool)
    .await?;
    let mut dockets_itemized = vec![];

    for docket in dockets {
        if !docket.docket_type.is_empty() {
            dockets_itemized.push((docket.uuid, docket.docket_type));
        }
    }

    let docket_update_closure = async |docket_item: (Uuid, String)| {
        let (docket_uuid, docket_type) = docket_item;
        // Split the docket_type on "-"
        let parts: Vec<&str> = docket_type.split('-').collect();

        if parts.len() >= 2 {
            let new_docket_type = parts[0].trim();
            let new_docket_subtype = parts[1].trim();

            // Update the docket with the new values
            let _ = sqlx::query!(
                "UPDATE dockets SET docket_type = $1, docket_subtype = $2 WHERE uuid = $3",
                new_docket_type,
                new_docket_subtype,
                docket_uuid
            )
            .execute(&pool)
            .await;
        }
    };
    let _iterators = stream::iter(dockets_itemized.into_iter())
        .map(docket_update_closure)
        .buffer_unordered(5)
        .count()
        .await;

    Ok(())
}
