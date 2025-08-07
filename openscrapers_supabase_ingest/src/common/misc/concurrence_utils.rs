use std::collections::BTreeMap;

use tokio::sync::{Mutex, Semaphore};
async fn blaahtest() {
    let mut test: BTreeMap<u64, u64> = BTreeMap::new();
    test.insert(0, 3);
}

async fn async_process_table<In, Out>(
    tasks: Vec<In>,
    process_individual: impl AsyncFn(In) -> Out,
    log_inorder_chunk_out: impl AsyncFn(Vec<Out>) -> (),
    simultaneous_tasks: usize,
) -> () {
    let concurrent_tasks = Semaphore::new(simultaneous_tasks);
    let mutex_log: Mutex<BTreeMap<usize, Out>> = Mutex::new(BTreeMap::new());

    let wrapped_process_individual =
        async |index: usize, input: In, log_ref: &Mutex<BTreeMap<usize, Out>>| {
            let output = process_individual(input).await;
            let mut log_lock = log_ref.lock().await;
            log_lock.insert(index, output);
            drop(log_lock);
        };
    // I want you to go ahead and continuously
}
