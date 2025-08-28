use async_trait::async_trait;
use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;
use sqlx::{PgConnection, PgTransaction, postgres::PgPoolOptions};
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

    let mut tx = pool.begin().await?;

    sqlx::query("SET LOCAL statement_timeout = 0;")
        .execute(&mut *tx)
        .await?;

    info!("Dropping existing tables");
    drop_existing_schema(&mut tx).await?;

    info!("Creating tables");
    create_schema(&mut tx).await?;

    tx.commit().await?;
    info!("Successfully recreated schema");

    Ok(())
}

pub async fn drop_existing_schema(tx: &mut PgConnection) -> anyhow::Result<()> {
    sqlx::query("DROP TABLE IF EXISTS public.docket_petitioned_by_org CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query("DROP TABLE IF EXISTS public.fillings_filed_by_org_relation CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query("DROP TABLE IF EXISTS public.fillings_on_behalf_of_org_relation CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query("DROP TABLE IF EXISTS public.attachments CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query("DROP TABLE IF EXISTS public.fillings CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query("DROP TABLE IF EXISTS public.dockets CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query("DROP TABLE IF EXISTS public.organizations CASCADE;")
        .execute(&mut *tx)
        .await?;
    Ok(())
}

pub async fn create_schema(tx: &mut PgConnection) -> anyhow::Result<()> {
    sqlx::query(
        "CREATE TABLE public.organizations (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          name character varying NOT NULL,
          aliases ARRAY NOT NULL,
          description character varying,
          artifical_person_type character varying,
          org_suffix character varying,
          CONSTRAINT organizations_pkey PRIMARY KEY (uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;

    sqlx::query(
        "CREATE TABLE public.dockets (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          docket_govid character varying NOT NULL DEFAULT ''::character varying UNIQUE,
          docket_subtype character varying,
          docket_description character varying,
          docket_title character varying,
          industry character varying,
          hearing_officer character varying,
          opened_date date NOT NULL,
          closed_date date,
          current_status character varying,
          assigned_judge character varying,
          docket_type character varying,
          petitioner_strings ARRAY,
          CONSTRAINT dockets_pkey PRIMARY KEY (uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;

    sqlx::query(
        "CREATE TABLE public.fillings (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          docket_uuid uuid NOT NULL,
          docket_govid character varying NOT NULL,
          individual_author_strings ARRAY NOT NULL,
          organization_author_strings ARRAY NOT NULL,
          filed_date date NOT NULL,
          filling_type character varying,
          filling_name character varying,
          filling_description character varying,
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          filling_govid character varying,
          CONSTRAINT fillings_pkey PRIMARY KEY (uuid),
          CONSTRAINT fillings_docket_uuid_fkey FOREIGN KEY (docket_uuid) REFERENCES public.dockets(uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;

    sqlx::query(
        "CREATE TABLE public.attachments (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          blake2b_hash character varying,
          parent_filling_uuid uuid NOT NULL,
          attachment_file_extension character varying,
          attachment_file_name character varying,
          attachment_title character varying,
          attachment_type character varying,
          attachment_subtype character varying,
          attachment_url character varying,
          openscrapers_id character varying NOT NULL DEFAULT ''::character varying UNIQUE,
          CONSTRAINT attachments_pkey PRIMARY KEY (uuid),
          CONSTRAINT attachments_parent_filling_uuid_fkey FOREIGN KEY (parent_filling_uuid) REFERENCES public.fillings(uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;

    sqlx::query(
        "CREATE TABLE public.docket_petitioned_by_org (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          docket_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          petitioner_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          CONSTRAINT docket_petitioned_by_org_pkey PRIMARY KEY (uuid),
          CONSTRAINT docket_petitioned_by_org_petitioner_uuid_fkey FOREIGN KEY (petitioner_uuid) REFERENCES public.organizations(uuid),
          CONSTRAINT docket_petitioned_by_org_docket_uuid_fkey FOREIGN KEY (docket_uuid) REFERENCES public.dockets(uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;

    sqlx::query(
        "CREATE TABLE public.fillings_filed_by_org_relation (
          relation_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          author_individual_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          filling_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          CONSTRAINT fillings_filed_by_org_relation_pkey PRIMARY KEY (relation_uuid),
          CONSTRAINT fillings_individual_authors_relation_filling_uuid_fkey FOREIGN KEY (filling_uuid) REFERENCES public.fillings(uuid),
          CONSTRAINT fillings_individual_authors_relatio_author_individual_uuid_fkey FOREIGN KEY (author_individual_uuid) REFERENCES public.organizations(uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;

    sqlx::query(
        "CREATE TABLE public.fillings_on_behalf_of_org_relation (
          relation_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          filling_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          author_organization_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          CONSTRAINT fillings_on_behalf_of_org_relation_pkey PRIMARY KEY (relation_uuid),
          CONSTRAINT fillings_organization_authors_rel_author_organization_uuid_fkey FOREIGN KEY (author_organization_uuid) REFERENCES public.organizations(uuid),
          CONSTRAINT fillings_organization_authors_relation_filling_uuid_fkey FOREIGN KEY (filling_uuid) REFERENCES public.fillings(uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;
    Ok(())
}
