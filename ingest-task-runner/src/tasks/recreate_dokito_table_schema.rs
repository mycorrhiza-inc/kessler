use async_trait::async_trait;
use schemars::JsonSchema;
use serde::Deserialize;
use serde_json::Value;
use sqlx::{PgConnection, postgres::PgPoolOptions};
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

    sqlx::query!("SET LOCAL statement_timeout = 0;")
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
    sqlx::query!("DROP TABLE IF EXISTS public.docket_petitioned_by_org CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query!("DROP TABLE IF EXISTS public.fillings_filed_by_org_relation CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query!("DROP TABLE IF EXISTS public.fillings_on_behalf_of_org_relation CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query!("DROP TABLE IF EXISTS public.attachments CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query!("DROP TABLE IF EXISTS public.fillings CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query!("DROP TABLE IF EXISTS public.dockets CASCADE;")
        .execute(&mut *tx)
        .await?;
    sqlx::query!("DROP TABLE IF EXISTS public.organizations CASCADE;")
        .execute(&mut *tx)
        .await?;
    Ok(())
}

pub async fn create_schema(tx: &mut PgConnection) -> anyhow::Result<()> {
    // organizations
    sqlx::query!(
        "CREATE TABLE public.organizations (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          name TEXT NOT NULL DEFAULT '',
          aliases TEXT[] NOT NULL DEFAULT '{}',
          description TEXT NOT NULL DEFAULT '',
          artifical_person_type TEXT NOT NULL DEFAULT '',
          org_suffix TEXT NOT NULL DEFAULT '',
          CONSTRAINT organizations_pkey PRIMARY KEY (uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;
    sqlx::query!("ALTER TABLE public.organizations ENABLE ROW LEVEL SECURITY;")
        .execute(&mut *tx)
        .await?;

    // dockets
    sqlx::query!(
        "CREATE TABLE public.dockets (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          docket_govid TEXT NOT NULL DEFAULT '' UNIQUE,
          docket_subtype TEXT NOT NULL DEFAULT '',
          docket_description TEXT NOT NULL DEFAULT '',
          docket_title TEXT NOT NULL DEFAULT '',
          industry TEXT NOT NULL DEFAULT '',
          hearing_officer TEXT NOT NULL DEFAULT '',
          opened_date date NOT NULL,
          closed_date date,
          current_status TEXT NOT NULL DEFAULT '',
          assigned_judge TEXT NOT NULL DEFAULT '',
          docket_type TEXT NOT NULL DEFAULT '',
          petitioner_strings TEXT[] NOT NULL DEFAULT '{}',
          CONSTRAINT dockets_pkey PRIMARY KEY (uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;
    sqlx::query!("ALTER TABLE public.dockets ENABLE ROW LEVEL SECURITY;")
        .execute(&mut *tx)
        .await?;

    // fillings
    sqlx::query!(
        "CREATE TABLE public.fillings (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          docket_uuid uuid NOT NULL,
          docket_govid TEXT NOT NULL DEFAULT '',
          individual_author_strings TEXT[] NOT NULL DEFAULT '{}',
          organization_author_strings TEXT[] NOT NULL DEFAULT '{}',
          filed_date date NOT NULL,
          filling_type TEXT NOT NULL DEFAULT '',
          filling_name TEXT NOT NULL DEFAULT '',
          filling_description TEXT NOT NULL DEFAULT '',
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          filling_govid TEXT NOT NULL DEFAULT '',
          CONSTRAINT fillings_pkey PRIMARY KEY (uuid),
          CONSTRAINT fillings_docket_uuid_fkey FOREIGN KEY (docket_uuid) REFERENCES public.dockets(uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;
    sqlx::query!("ALTER TABLE public.fillings ENABLE ROW LEVEL SECURITY;")
        .execute(&mut *tx)
        .await?;

    // attachments
    sqlx::query!(
        "CREATE TABLE public.attachments (
          uuid uuid NOT NULL DEFAULT gen_random_uuid(),
          created_at timestamp with time zone NOT NULL DEFAULT now(),
          updated_at timestamp with time zone NOT NULL DEFAULT now(),
          blake2b_hash TEXT NOT NULL DEFAULT '',
          parent_filling_uuid uuid NOT NULL,
          attachment_file_extension TEXT NOT NULL DEFAULT '',
          attachment_file_name TEXT NOT NULL DEFAULT '',
          attachment_title TEXT NOT NULL DEFAULT '',
          attachment_type TEXT NOT NULL DEFAULT '',
          attachment_subtype TEXT NOT NULL DEFAULT '',
          attachment_url TEXT NOT NULL DEFAULT '',
          openscrapers_id TEXT NOT NULL DEFAULT '' UNIQUE,
          CONSTRAINT attachments_pkey PRIMARY KEY (uuid),
          CONSTRAINT attachments_parent_filling_uuid_fkey FOREIGN KEY (parent_filling_uuid) REFERENCES public.fillings(uuid)
        );",
    )
    .execute(&mut *tx)
    .await?;
    sqlx::query!("ALTER TABLE public.attachments ENABLE ROW LEVEL SECURITY;")
        .execute(&mut *tx)
        .await?;

    // docket_petitioned_by_org
    sqlx::query!(
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
    sqlx::query!("ALTER TABLE public.docket_petitioned_by_org ENABLE ROW LEVEL SECURITY;")
        .execute(&mut *tx)
        .await?;

    // fillings_filed_by_org_relation
    sqlx::query!(
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
    sqlx::query!("ALTER TABLE public.fillings_filed_by_org_relation ENABLE ROW LEVEL SECURITY;")
        .execute(&mut *tx)
        .await?;

    // fillings_on_behalf_of_org_relation
    sqlx::query!(
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
    sqlx::query!(
        "ALTER TABLE public.fillings_on_behalf_of_org_relation ENABLE ROW LEVEL SECURITY;"
    )
    .execute(&mut *tx)
    .await?;

    Ok(())
}
