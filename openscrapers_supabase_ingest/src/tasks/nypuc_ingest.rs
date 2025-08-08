use async_trait::async_trait;
use serde_json::Value;

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
    Ok(())
}

pub async fn ingest_nypuc_case(case: GenericCase) -> anyhow::Result<()> {
    // this will take in a case from openscrapers of a schema GenericCase. Ingest it into a
    // postgres database with schema in
    // /home/nicole/Documents/mycorrhiza/kessler/openscrapers_supabase_ingest/supabase_sql_schema.sql
    // use the sqlx library for everything.
    //
    // This should mainly get broken down into a bunch of different steps.
    // Go ahead and create the docket object, (check and see if there is a govid match in the
    // database, if so do a cascade delete.)]
    //
    // Create a bunch of filling objects each referencing back to the docket uuid.
    //
    // Then create a bunch of attachment objects that each reference back to that fileID
    //
    // Then check the artifical persons database to see if any entry has an organizationName that
    // matches the name column in artifical persons, if not go ahead and create a new artifical
    // person. And link that up to the filling via the relations table.
    // Do the same for the individual authors.
    todo!()
}
