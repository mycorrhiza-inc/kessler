use std::{env, sync::LazyLock};

use async_trait::async_trait;
use axum::Json;
use futures::stream::{self, StreamExt};
use reqwest::Client;
use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;
use sqlx::{PgPool, Pool, Postgres, postgres::PgPoolOptions, types::Uuid};

use crate::{
    common::{
        llm_deepinfra::org_split_from_dump,
        misc::map_empty,
        tasks::{
            ExecuteUserTask, TaskStatusDisplay, routing::PriorityExtractor,
            workers::add_task_to_queue,
        },
    },
    types::openscrapers::RawGenericCase,
};
use tracing::info;

#[derive(Clone, Copy, Default, Deserialize, JsonSchema)]
pub struct NyPucIngestFull {}

#[async_trait]
impl ExecuteUserTask for NyPucIngestFull {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = get_all_ny_puc_data().await;
        match res {
            Ok(()) => {
                info!("Nypuc ingest completed.");
                Ok("Task Completed Successfully".into())
            }
            Err(err) => {
                tracing::error!(error= % err, error_debug= ?err,"Encountered error in ny_ingest");
                Err(err.to_string().into())
            }
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

pub async fn add_nypuc_all_task(
    Json(PriorityExtractor { priority }): Json<PriorityExtractor>,
) -> Json<TaskStatusDisplay> {
    let ny_ingest = NyPucIngestFull::default();
    let taskinfo = add_task_to_queue(ny_ingest, priority).await;
    Json(taskinfo.into())
}

pub async fn get_all_ny_puc_data() -> anyhow::Result<()> {
    info!("Got request to ingest all nypuc data.");
    let reqwest_client = Client::new();

    // Drop all existing tables first
    delete_all_data().await?;
    info!("Successfully deleted all old case data.");

    // Get the list of case IDs
    let case_ids: Vec<String> = reqwest_client
        .get("http://localhost:33399/public/caselist/ny/ny_puc/all")
        .send()
        .await?
        .json()
        .await?;
    info!(length=%case_ids.len(),"Got list of all cases");

    let db_url = &**DEFAULT_POSTGRES_CONNECTION_URL;
    // let options =
    let pool = PgPoolOptions::new()
        .max_connections(60)
        .connect(db_url)
        .await?;

    info!("Created pg pool");

    let execute_case_wraped = async |case_id| {
        let client = reqwest_client.clone();
        let url = format!("http://localhost:33399/public/cases/ny/ny_puc/{case_id}");
        let res = client.get(&url).send().await;

        match res {
            Ok(response) => {
                let response_bytes = response
                    .text()
                    .await
                    .unwrap_or("encountered error getting raw response bytes".to_string());
                let case_res = serde_json::from_str::<RawGenericCase>(&response_bytes);
                match case_res {
                    Ok(case) => {
                        if let Err(e) = ingest_nypuc_case(case, &pool).await {
                            tracing::error!(case_id = %case_id, error = %e, error_debug = ?e, "Failed to ingest case");
                        }
                    }
                    Err(e) => {
                        tracing::error!(case_id = %case_id, error = %e, error_debug = ?e, raw_response =%response_bytes[0..400],"Failed to parse case")
                    }
                }
            }
            Err(e) => {
                tracing::error!(url,case_id = %case_id, error = %e, error_debug = ?e,"Failed to fetch case")
            }
        }
    };

    // Create a stream of futures to fetch and ingest each case concurrently
    let futures_count = stream::iter(case_ids)
        .map(execute_case_wraped)
        .buffer_unordered(10)
        .count()
        .await;
    info!(futures_count, "Successfully completed all futures.");
    Ok(())
}

static DEFAULT_POSTGRES_CONNECTION_URL: LazyLock<String> = LazyLock::new(|| {
    env::var("POSTGRES_CONNECTION")
        .or(env::var("DATABASE_URL"))
        .expect("POSTGRES_CONNECTION or DATABASE_URL should be set.")
});

pub async fn ingest_nypuc_case(case: RawGenericCase, pool: &Pool<Postgres>) -> anyhow::Result<()> {
    // Check for existing docket and delete if found
    let existing_docket: Option<Uuid> = sqlx::query_scalar!(
        "SELECT uuid FROM dockets WHERE docket_govid = $1",
        &case.case_govid.as_str()
    )
    .fetch_optional(pool)
    .await?;

    if let Some(docket_uuid) = existing_docket {
        sqlx::query!("DELETE FROM dockets WHERE uuid = $1", docket_uuid)
            .execute(pool)
            .await?;
    }
    // FIXME: Delete all this shit once its actually processed once, this should be an openscraper
    // responsibility.
    let petitioner_str = map_empty(&*case.petitioner);
    let mut petitioner_list = vec![];
    if let Some(petitioner_str) = petitioner_str {
        petitioner_list = org_split_from_dump(petitioner_str).await?;
    }

    // Create new docket
    let docket_uuid: Uuid = sqlx::query_scalar!(
        "INSERT INTO dockets (docket_govid, docket_description, docket_title, industry, hearing_officer, opened_date, closed_date, petitioner_strings, docket_type, docket_subtype )
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING uuid",
        &case.case_govid.as_str(),
        map_empty(&case.description),
        &case.case_name,
        map_empty(&case.industry),
        map_empty(&case.hearing_officer),
        case.opened_date,
        case.closed_date,
        &petitioner_list,
        map_empty(&case.case_type),
        Option::<String>::None
    )
    .fetch_one(pool)
    .await?;
    for petitioner in petitioner_list.iter() {
        let petitioner_uuid = fetch_or_insert_new_orgstring(petitioner, pool).await?;
        sqlx::query!(
            "INSERT INTO docket_petitioned_by_org (docket_uuid, petitioner_uuid) VALUES ($1,$2)",
            docket_uuid,
            petitioner_uuid
        )
        .execute(pool)
        .await?;
    }

    for filling in case.filings {
        let filling_uuid: Uuid = sqlx::query_scalar!(
            "INSERT INTO fillings (docket_uuid, docket_govid, individual_author_strings, organization_author_strings, filed_date, filling_type, filling_name, filling_description)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uuid",
            docket_uuid,
            &case.case_govid.as_str(),
            &filling.individual_authors.iter().map(|s| s.to_string()).collect::<Vec<_>>(),
            &filling.organization_authors.iter().map(|s| s.to_string()).collect::<Vec<_>>(),
            filling.filed_date,
            &filling.filing_type,
            &filling.name,
            map_empty(&filling.description),
        )
        .fetch_one(pool)
        .await?;

        for attachment in filling.attachments {
            if let Some(hash) = attachment.hash {
                sqlx::query!(
                "INSERT INTO attachments (parent_filling_uuid, blake2b_hash, attachment_file_extension, attachment_file_name, attachment_title, attachment_url, created_at, updated_at)
                 VALUES ($1, $2, $3, $4, $5, $6, now(), now())",
                filling_uuid,
                hash.to_string(),
                &attachment.document_extension.to_string(),
                &attachment.name,
                &attachment.name,
                &attachment.url
            ).execute(pool)
            .await?;
            } else {
                // tracing::warn!(
                //     ?attachment,
                //     "Encountered attachment with missing hash, inserting regardless"
                // );
                let nullstr: Option<&str> = None;
                sqlx::query!(
                "INSERT INTO attachments (parent_filling_uuid, blake2b_hash, attachment_file_extension, attachment_file_name, attachment_title, attachment_url, created_at, updated_at)
                 VALUES ($1, $2, $3, $4, $5, $6, now(), now())",
                filling_uuid,
                nullstr,
                &attachment.document_extension.to_string(),
                &attachment.name,
                &attachment.name,
                &attachment.url
            ).execute(pool).await?;
            }
        }

        for indiv_author in filling.individual_authors {
            let org_uuid = fetch_or_insert_new_orgstring(indiv_author.as_str(), pool).await?;

            sqlx::query!(
                "INSERT INTO fillings_filed_by_org_relation (author_individual_uuid, filling_uuid) VALUES ($1, $2)",
                org_uuid,
                filling_uuid
            )
            .execute(pool)
            .await?;
        }

        for org_author in filling.organization_authors {
            let org_uuid = fetch_or_insert_new_orgstring(org_author.as_str(), pool).await?;
            sqlx::query!(
                "INSERT INTO fillings_on_behalf_of_org_relation (author_organization_uuid, filling_uuid) VALUES ($1, $2)",
                org_uuid,
                filling_uuid
            )
            .execute(pool)
            .await?;
        }
    }
    tracing::info!(govid=%case.case_govid, uuid=%docket_uuid,"Successfully processed case with no errors");

    Ok(())
}

async fn fetch_or_insert_new_orgstring(
    org_author: &str,
    pool: &Pool<Postgres>,
) -> Result<Uuid, anyhow::Error> {
    let org_record: Option<Uuid> = sqlx::query_scalar!(
        "SELECT uuid FROM organizations WHERE name = $1 AND artifical_person_type = 'organization'",
        org_author
    )
    .fetch_optional(pool)
    .await?;

    let org_uuid = if let Some(org_record) = org_record {
        org_record
    } else {
        let new_org: Uuid = sqlx::query_scalar!(
                    "INSERT INTO organizations (name, artifical_person_type, aliases) VALUES ($1, 'organization', $2) RETURNING uuid",
                    org_author,
                    &vec![org_author.to_string()]
                )
                .fetch_one(pool)
                .await?;
        new_org
    };
    Ok(org_uuid)
}

pub async fn delete_all_data() -> anyhow::Result<()> {
    let db_url = &**DEFAULT_POSTGRES_CONNECTION_URL;
    let pool = PgPool::connect(db_url).await?;

    // Drop all data from tables in the correct order to avoid foreign key constraint violations
    // Start with the relation tables
    sqlx::query!("DELETE FROM fillings_filed_by_org_relation")
        .execute(&pool)
        .await?;

    sqlx::query!("DELETE FROM fillings_on_behalf_of_org_relation")
        .execute(&pool)
        .await?;

    // Then attachments
    sqlx::query!("DELETE FROM attachments")
        .execute(&pool)
        .await?;

    // Then fillings
    sqlx::query!("DELETE FROM fillings").execute(&pool).await?;

    // Then dockets
    sqlx::query!("DELETE FROM dockets").execute(&pool).await?;

    // Finally organizations
    sqlx::query!("DELETE FROM organizations")
        .execute(&pool)
        .await?;

    Ok(())
}
