use async_trait::async_trait;
use serde_json::Value;

use crate::common::task_workers::ExecuteUserTask;

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
    todo!()
}
