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
        misc::{fmap_empty, map_empty},
        tasks::{
            ExecuteUserTask, TaskStatusDisplay, routing::PriorityExtractor,
            workers::add_task_to_queue,
        },
    },
    types::openscrapers::GenericCase,
};

#[derive(Clone, Copy, Default, Deserialize, JsonSchema)]
pub struct NyPucIngestFull {}

#[async_trait]
impl ExecuteUserTask for NyPucIngestFull {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = get_all_ny_puc_data().await;
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

pub async fn add_nypuc_all_task(
    Json(PriorityExtractor { priority }): Json<PriorityExtractor>,
) -> Json<TaskStatusDisplay> {
    let ny_ingest = NyPucIngestFull::default();
    let taskinfo = add_task_to_queue(ny_ingest, priority).await;
    Json(taskinfo.into())
}

pub async fn get_all_ny_puc_data() -> anyhow::Result<()> {
    let reqwest_client = Client::new();

    // Drop all existing tables first
    delete_all_data().await?;

    // Get the list of case IDs
    let case_ids: Vec<String> = reqwest_client
        .get("http://localhost:33399/public/caselist/ny/ny_puc/all")
        .send()
        .await?
        .json()
        .await?;

    let db_url = &**DEFAULT_POSTGRES_CONNECTION_URL;
    // let options =
    let pool = PgPoolOptions::new()
        .max_connections(30)
        .connect(db_url)
        .await?;

    // Create a stream of futures to fetch and ingest each case concurrently
    let futures = stream::iter(case_ids)
        .map(async |case_id| {
            let client = reqwest_client.clone();
            let url = format!("http://localhost:33399/public/cases/ny/ny_puc/{case_id}");
            let res = client.get(&url).send().await;

            match res {
                Ok(response) => {
                    let case_res = response.json::<GenericCase>().await;
                    match case_res {
                        Ok(case) => {
                            if let Err(e) = ingest_nypuc_case(case,&pool).await {
                                tracing::error!(case_id = %case_id, error = %e, "Failed to ingest case");
                            }
                        }
                        Err(e) => {
                            tracing::error!(case_id = %case_id, error = %e, "Failed to parse case")
                        }
                    }
                }
                Err(e) => tracing::error!(case_id = %case_id, error = %e, "Failed to fetch case"),
            }
        })
        .buffer_unordered(15); // Process up to 10 requests concurrently

    // Wait for all futures to complete
    futures.for_each(|_| async {}).await;

    Ok(())
}

static DEFAULT_POSTGRES_CONNECTION_URL: LazyLock<String> = LazyLock::new(|| {
    env::var("POSTGRES_CONNECTION")
        .or(env::var("DATABASE_URL"))
        .expect("POSTGRES_CONNECTION or DATABASE_URL should be set.")
});

pub async fn ingest_nypuc_case(case: GenericCase, pool: &Pool<Postgres>) -> anyhow::Result<()> {
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

    // Create new docket
    let docket_uuid: Uuid = sqlx::query_scalar!(
        "INSERT INTO dockets (docket_govid, docket_description, docket_title, industry, petitioner, hearing_officer, opened_date, closed_date)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uuid",
        &case.case_govid.as_str(),
        fmap_empty(Some(&case.description)),
        &case.case_name,
        fmap_empty(Some(&case.industry)),
        fmap_empty(Some(&case.petitioner)),
        fmap_empty(Some(&case.hearing_officer)),
        case.opened_date,
        case.closed_date
    )
    .fetch_one(pool)
    .await?;

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
                tracing::warn!(
                    ?attachment,
                    "Encountered attachment with missing hash, inserting regardless"
                );
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

        for author in filling.individual_authors {
            let person: Option<Uuid> = sqlx::query_scalar!(
                "SELECT uuid FROM artifical_persons WHERE name = $1 AND is_human = true",
                &author.as_str()
            )
            .fetch_optional(pool)
            .await?;

            let person_uuid = if let Some(person) = person {
                person
            } else {
                let new_person: Uuid = sqlx::query_scalar!(
                    "INSERT INTO artifical_persons (name, is_human, is_corporate_entity, aliases) VALUES ($1, true, false, $2) RETURNING uuid",
                    &author.as_str(),
                    &vec![author.to_string()]
                )
                .fetch_one(pool)
                .await?;
                new_person
            };

            sqlx::query!(
                "INSERT INTO fillings_individual_authors_relation (author_individual_uuid, filling_uuid) VALUES ($1, $2)",
                person_uuid,
                filling_uuid
            )
            .execute(pool)
            .await?;
        }

        for org in filling.organization_authors {
            let person: Option<Uuid> = sqlx::query_scalar!(
                "SELECT uuid FROM artifical_persons WHERE name = $1 AND is_corporate_entity = true",
                &org.as_str()
            )
            .fetch_optional(pool)
            .await?;

            let person_uuid = if let Some(person) = person {
                person
            } else {
                let new_person: Uuid = sqlx::query_scalar!(
                    "INSERT INTO artifical_persons (name, is_human, is_corporate_entity, aliases) VALUES ($1, false, true, $2) RETURNING uuid",
                    &org.as_str(),
                    &vec![org.to_string()]
                )
                .fetch_one(pool)
                .await?;
                new_person
            };

            sqlx::query!(
                "INSERT INTO fillings_organization_authors_relation (author_organization_uuid, filling_uuid) VALUES ($1, $2)",
                person_uuid,
                filling_uuid
            )
            .execute(pool)
            .await?;
        }
    }

    Ok(())
}

pub async fn delete_all_data() -> anyhow::Result<()> {
    let db_url = &**DEFAULT_POSTGRES_CONNECTION_URL;
    let pool = PgPool::connect(db_url).await?;

    // Drop all data from tables in the correct order to avoid foreign key constraint violations
    // Start with the relation tables
    sqlx::query!("DELETE FROM fillings_individual_authors_relation")
        .execute(&pool)
        .await?;

    sqlx::query!("DELETE FROM fillings_organization_authors_relation")
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

    // Finally artificial persons
    sqlx::query!("DELETE FROM artifical_persons")
        .execute(&pool)
        .await?;

    Ok(())
}
