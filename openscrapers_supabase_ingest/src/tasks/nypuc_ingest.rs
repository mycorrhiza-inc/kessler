use std::env;

use async_trait::async_trait;
use serde_json::Value;
use sqlx::{PgPool, types::Uuid};

use crate::{common::task_workers::ExecuteUserTask, types::openscrapers::GenericCase};

struct NyPucIngestFull {}

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
}

pub async fn get_all_ny_puc_data() -> anyhow::Result<()> {
    // This ny_puc ingest function will get a bunch of lists from the openscrapers api
    // Dont implement this for now
    // if you run this curl command
    // curl -X 'GET' \
    // 'http://localhost:33399/public/caselist/ny/ny_puc/all' \
    // -H 'accept: application/json'
    // you get a response like this
    // [
    // "10-00036",
    // "10-00320",
    // "10-00529",
    // "10-01102",
    // "10-02623",
    // "10-M-0082",
    // "10-M-0186",
    // "10-M-0218",
    // "10-M-0365",
    // "10-T-0019",
    // "10-T-0453",
    // "11-00714",
    // "11-00751",]
    // thats a bunch of govids for the dockets,
    // then if you hit this endpoint:
    // curl -X 'GET' \
    // 'http://localhost:33399/public/cases/ny/ny_puc/10-00036' \
    // -H 'accept: application/json'
    // You get a GenericCaseLegacy type. Convert that to a GenericCase, and run the ingest on it.
    // Do this for all case results, and maybe add some async so its faster.
    Ok(())
}

pub async fn ingest_nypuc_case(case: GenericCase) -> anyhow::Result<()> {
    let db_url = env::var("DATABASE_URL")?;
    let pool = PgPool::connect(&db_url).await?;

    // Check for existing docket and delete if found
    let existing_docket: Option<(Uuid,)> =
        sqlx::query_as("SELECT uuid FROM dockets WHERE docket_govid = $1")
            .bind(&case.case_number)
            .fetch_optional(&pool)
            .await?;

    if let Some((docket_uuid,)) = existing_docket {
        sqlx::query("DELETE FROM dockets WHERE uuid = $1")
            .bind(docket_uuid)
            .execute(&pool)
            .await?;
    }

    // Create new docket
    let docket_uuid: (Uuid,) = sqlx::query_as(
        "INSERT INTO dockets (docket_govid, docket_description, docket_title, industry, petitioner, hearing_officer, opened_date, closed_date)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uuid"
    )
    .bind(&case.case_number)
    .bind(&case.description)
    .bind(&case.case_name)
    .bind(&case.industry)
    .bind(&case.petitioner)
    .bind(&case.hearing_officer)
    .bind(case.opened_date)
    .bind(case.closed_date)
    .fetch_one(&pool)
    .await?;

    for filling in case.filings {
        let filling_uuid: (Uuid,) = sqlx::query_as(
            "INSERT INTO fillings (docket_uuid, docket_govid, individual_author_strings, organization_author_strings, filed_date, filling_type, filling_name, filling_description)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING uuid"
        )
        .bind(docket_uuid.0)
        .bind(&case.case_number)
        .bind(&filling.individual_authors)
        .bind(&filling.organization_authors)
        .bind(filling.filed_date)
        .bind(&filling.filing_type)
        .bind(&filling.name)
        .bind(&filling.description)
        .fetch_one(&pool)
        .await?;

        for attachment in filling.attachments {
            if let Some(hash) = attachment.hash {
                sqlx::query(
                "INSERT INTO attachments (parent_filling_uuid, blake2b_hash, attachment_file_extension, attachment_file_name, attachment_title, attachment_url)
                 VALUES ($1, $2, $3, $4, $5, $6)"
            ).bind(filling_uuid.0)
            .bind(hash.to_string())
            .bind(&attachment.document_extension)
            .bind(&attachment.name)
            .bind(&attachment.name)
            .bind(&attachment.url)
            .execute(&pool)
            .await?;
            }
        }

        for author in filling.individual_authors {
            let person: Option<(Uuid,)> = sqlx::query_as(
                "SELECT uuid FROM artifical_persons WHERE name = $1 AND is_human = true",
            )
            .bind(&author)
            .fetch_optional(&pool)
            .await?;

            let person_uuid = if let Some(person) = person {
                person.0
            } else {
                let new_person: (Uuid,) = sqlx::query_as(
                    "INSERT INTO artifical_persons (name, is_human, is_corporate_entity) VALUES ($1, true, false) RETURNING uuid"
                )
                .bind(&author)
                .fetch_one(&pool)
                .await?;
                new_person.0
            };

            sqlx::query(
                "INSERT INTO fillings_individual_authors_relation (author_individual_uuid, filling_uuid) VALUES ($1, $2)"
            )
            .bind(person_uuid)
            .bind(filling_uuid.0)
            .execute(&pool)
            .await?;
        }

        for org in filling.organization_authors {
            let person: Option<(Uuid,)> = sqlx::query_as(
                "SELECT uuid FROM artifical_persons WHERE name = $1 AND is_corporate_entity = true",
            )
            .bind(&org)
            .fetch_optional(&pool)
            .await?;

            let person_uuid = if let Some(person) = person {
                person.0
            } else {
                let new_person: (Uuid,) = sqlx::query_as(
                    "INSERT INTO artifical_persons (name, is_human, is_corporate_entity) VALUES ($1, false, true) RETURNING uuid"
                )
                .bind(&org)
                .fetch_one(&pool)
                .await?;
                new_person.0
            };

            sqlx::query(
                "INSERT INTO fillings_organization_authors_relation (author_organization_uuid, filling_uuid) VALUES ($1, $2)"
            )
            .bind(person_uuid)
            .bind(filling_uuid.0)
            .execute(&pool)
            .await?;
        }
    }

    Ok(())
}
