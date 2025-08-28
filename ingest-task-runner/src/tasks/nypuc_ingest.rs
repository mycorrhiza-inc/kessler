use std::{env, sync::LazyLock};

use async_trait::async_trait;
use futures::stream::{self, StreamExt};
use rand::{SeedableRng, rngs::SmallRng, seq::SliceRandom};
use reqwest::Client;
use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;
use sqlx::{PgPool, Pool, Postgres, postgres::PgPoolOptions, types::Uuid};

use mycorrhiza_common::{misc::map_empty, tasks::ExecuteUserTask};
use tracing::{info, warn};

use crate::types::openscrapers::{OrgName, ProcessedGenericDocket};

#[derive(Clone, Copy, Default, Deserialize, JsonSchema)]
pub struct NyPucIngestPurgePrevious {}

#[async_trait]
impl ExecuteUserTask for NyPucIngestPurgePrevious {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = get_all_ny_puc_data(true).await;
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
        "ingest_nypuc_purge_previous"
    }
    fn get_task_label_static() -> &'static str
    where
        Self: Sized,
    {
        "ingest_nypuc_purge_previous"
    }
}

#[derive(Clone, Copy, Default, Deserialize, JsonSchema)]
pub struct NyPucIngestGetMissingDockets {}
#[async_trait]
impl ExecuteUserTask for NyPucIngestGetMissingDockets {
    async fn execute_task(self: Box<Self>) -> Result<Value, Value> {
        let res = get_all_ny_puc_data(false).await;
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
        "ingest_nypuc_get_missing_dockets"
    }
    fn get_task_label_static() -> &'static str
    where
        Self: Sized,
    {
        "ingest_nypuc_get_missing_dockets"
    }
}

pub async fn get_all_ny_puc_data(purge_data: bool) -> anyhow::Result<()> {
    info!("Got request to ingest all nypuc data.");
    let reqwest_client = Client::new();

    let db_url = &**DEFAULT_POSTGRES_CONNECTION_URL;
    let pool = PgPoolOptions::new()
        .max_connections(40)
        .connect(db_url)
        .await?;
    info!("Created pg pool");

    // Drop all existing tables first
    if purge_data {
        delete_all_data(&pool).await?;
        info!("Successfully deleted all old case data.");
    }

    // Get the list of case IDs
    let mut case_ids: Vec<String> = reqwest_client
        .get("http://localhost:33399/public/caselist/ny/ny_puc/all")
        .send()
        .await?
        .json()
        .await?;
    info!(length=%case_ids.len(),"Got list of all cases");

    // randomize list before actual ingest process.
    let mut rng = SmallRng::from_os_rng();
    case_ids.shuffle(&mut rng);

    // We can set this to always true since we just purged the dataset.
    let ignore_existing = true;

    let execute_case_wraped =
        async |case_id: String| ingest_wrapped_ny_data(&case_id, &pool, ignore_existing).await;

    // Create a stream of futures to fetch and ingest each case concurrently
    let futures_count = stream::iter(case_ids)
        .map(execute_case_wraped)
        .buffer_unordered(10)
        .count()
        .await;
    info!(futures_count, "Successfully completed all futures.");
    Ok(())
}

async fn ingest_wrapped_ny_data(case_id: &str, pool: &PgPool, ignore_existing: bool) {
    let client = Client::new();
    let url = format!("http://localhost:33399/public/cases/ny/ny_puc/{case_id}");
    let res = client.get(&url).send().await;

    match res {
        Ok(response) => {
            let response_bytes = response
                .text()
                .await
                .unwrap_or("encountered error getting raw response bytes".to_string());
            let case_res = serde_json::from_str::<ProcessedGenericDocket>(&response_bytes);
            match case_res {
                Ok(case) => {
                    const CASE_RETRIES: usize = 3;
                    if let Err(e) =
                        ingest_case_with_retries(&case, pool, ignore_existing, CASE_RETRIES).await
                    {
                        tracing::error!(case_id = %case_id, error = %e, error_debug = ?e, "Failed to ingest case, dispite retries.");
                    }
                }
                Err(e) => {
                    let subslice = if response_bytes.len() > 301 {
                        &response_bytes[..300]
                    } else {
                        &response_bytes
                    };
                    tracing::error!(case_id = %case_id, error = %e, error_debug = ?e, raw_response =%subslice,"Failed to parse case")
                }
            }
        }
        Err(e) => {
            tracing::error!(url,case_id = %case_id, error = %e, error_debug = ?e,"Failed to fetch case")
        }
    }
}

pub static DEFAULT_POSTGRES_CONNECTION_URL: LazyLock<String> = LazyLock::new(|| {
    env::var("POSTGRES_CONNECTION")
        .or(env::var("DATABASE_URL"))
        .expect("POSTGRES_CONNECTION or DATABASE_URL should be set.")
});

pub async fn ingest_case_with_retries(
    case: &ProcessedGenericDocket,
    pool: &Pool<Postgres>,
    ignore_existing: bool,
    tries: usize,
) -> anyhow::Result<()> {
    let mut return_res = Ok(());
    for remaining_tries in (0..tries).rev() {
        match ingest_nypuc_case(case, pool, ignore_existing).await {
            Ok(val) => return Ok(val),
            Err(err) => {
                warn!(docket_govid=%case.case_govid, %remaining_tries,"Encountered error while processing docket, retrying.");
                return_res = Err(err);
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
                    info!(%docket_uuid, %case.case_govid,"Successfully deleted corrupted case data");
                } else {
                    info!(%case.case_govid,"Case ingest errored, but could not find corrupt case data.");
                }
            }
        }
    }
    return_res
}

pub async fn ingest_nypuc_case(
    case: &ProcessedGenericDocket,
    pool: &Pool<Postgres>,
    ignore_existing: bool,
) -> anyhow::Result<()> {
    // Check for existing docket and delete if found
    let existing_docket: Option<Uuid> = sqlx::query_scalar!(
        "SELECT uuid FROM dockets WHERE docket_govid = $1",
        &case.case_govid.as_str()
    )
    .fetch_optional(pool)
    .await?;

    if let Some(docket_uuid) = existing_docket {
        if ignore_existing {
            return Ok(());
        } else {
            sqlx::query!("DELETE FROM dockets WHERE uuid = $1", docket_uuid)
                .execute(pool)
                .await?;
        }
    }
    // FIXME: Delete all this shit once its actually processed once, this should be an openscraper
    // responsibility.
    let petitioner_list: Vec<OrgName> = vec![];
    let petitioner_strings = petitioner_list
        .iter()
        .map(|n| n.name.to_string())
        .collect::<Vec<_>>();

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
        &petitioner_strings,
        map_empty(&case.case_type),
        Option::<String>::None
    )
    .fetch_one(pool)
    .await?;
    for petitioner in petitioner_list.iter() {
        let petitioner_uuid = fetch_or_insert_new_orgname(petitioner, pool).await?;
        sqlx::query!(
            "INSERT INTO docket_petitioned_by_org (docket_uuid, petitioner_uuid) VALUES ($1,$2)",
            docket_uuid,
            petitioner_uuid
        )
        .execute(pool)
        .await?;
    }

    for (_, filling) in case.filings.iter() {
        let individual_author_strings = filling
            .individual_authors
            .iter()
            .map(|s| s.name.to_string())
            .collect::<Vec<_>>();
        let organization_author_strings = filling
            .individual_authors
            .iter()
            .map(|s| s.name.to_string())
            .collect::<Vec<_>>();
        let filling_uuid: Uuid = sqlx::query_scalar!(
            "INSERT INTO fillings (docket_uuid, docket_govid, individual_author_strings, organization_author_strings, filed_date, filling_type, filling_name, filling_description)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uuid",
            docket_uuid,
            &case.case_govid.as_str(),
            &individual_author_strings,
            &organization_author_strings,
            filling.filed_date,
            &filling.filing_type,
            &filling.name,
            map_empty(&filling.description),
        )
        .fetch_one(pool)
        .await?;

        for (_, attachment) in filling.attachments.iter() {
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

        for indiv_author in filling.individual_authors.iter() {
            let org_uuid = fetch_or_insert_new_orgname(indiv_author, pool).await?;

            sqlx::query!(
                "INSERT INTO fillings_filed_by_org_relation (author_individual_uuid, filling_uuid) VALUES ($1, $2)",
                org_uuid,
                filling_uuid
            )
            .execute(pool)
            .await?;
        }

        for org_author in filling.organization_authors.iter() {
            let org_uuid = fetch_or_insert_new_orgname(org_author, pool).await?;
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

async fn fetch_or_insert_new_orgname(
    org_author: &OrgName,
    pool: &Pool<Postgres>,
) -> Result<Uuid, anyhow::Error> {
    let org_author_str = org_author.name.as_str();
    let org_record  = sqlx::query!(
        "SELECT uuid, org_suffix FROM organizations WHERE name = $1 AND artifical_person_type = 'organization'",
        org_author_str,
    ).fetch_optional(pool)
    .await?;

    let org_uuid = if let Some(org_record) = org_record {
        if org_record.org_suffix.is_none() && !org_author.suffix.is_empty() {
            let _ = sqlx::query!(
                "UPDATE organizations SET org_suffix = $1 WHERE uuid = $2",
                &org_author.suffix,
                &org_record.uuid
            )
            .execute(pool)
            .await?;
        };
        org_record.uuid
    } else {
        let org_suffix = map_empty(&*org_author.suffix);
        let new_org: Uuid = sqlx::query_scalar!(
                    "INSERT INTO organizations (name, artifical_person_type, aliases, org_suffix) VALUES ($1, 'organization', $2, $3) RETURNING uuid",
                    org_author_str,
                    &vec![org_author_str.to_string()],
                    org_suffix,
                )
                .fetch_one(pool)
                .await?;
        new_org
    };
    Ok(org_uuid)
}

pub async fn delete_all_data(pool: &PgPool) -> anyhow::Result<()> {
    info!("Starting full data deletion...");

    // Start a transaction
    let mut tx = pool.begin().await?;

    // Disable statement timeout just for this transaction
    sqlx::query("SET LOCAL statement_timeout = 0;")
        .execute(&mut *tx)
        .await?;
    info!("Disabled statement_timeout for this transaction");

    // Drop relation tables first
    info!("Deleting from fillings_filed_by_org_relation");
    sqlx::query!("TRUNCATE fillings_filed_by_org_relation")
        .execute(&mut *tx)
        .await?;

    info!("Deleting from fillings_on_behalf_of_org_relation");
    sqlx::query!("TRUNCATE fillings_on_behalf_of_org_relation")
        .execute(&mut *tx)
        .await?;

    // Attachments
    info!("Deleting from attachments");
    sqlx::query!("TRUNCATE attachments")
        .execute(&mut *tx)
        .await?;

    // Organizations
    info!("Deleting from organizations");
    sqlx::query!("TRUNCATE organizations")
        .execute(&mut *tx)
        .await?;

    // Fillings
    info!("Deleting from fillings");
    sqlx::query!("TRUNCATE fillings").execute(&mut *tx).await?;

    // Dockets
    info!("Deleting from dockets");
    sqlx::query!("TRUNCATE dockets").execute(&mut *tx).await?;

    // Commit once everything is successful
    tx.commit().await?;
    info!("All data deleted successfully ✅");

    Ok(())
}
