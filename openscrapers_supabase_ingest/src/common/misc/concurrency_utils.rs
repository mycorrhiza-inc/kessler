use std::{
    collections::{BTreeMap, VecDeque},
    ops::DerefMut,
    time::Duration,
};

use tokio::sync::Mutex;


async fn wrap_individual_process<In, Out>(
    input: (usize, In),
    log_ref: &Mutex<BTreeMap<usize, Out>>,
    process_individual: &impl AsyncFn(In) -> Out,
) {
    let (index, val) = input;
    let output = process_individual(val).await;
    let mut log_lock = log_ref.lock().await;
    log_lock.insert(index, output);
    drop(log_lock);
}

async fn log_out_from_log<Out>(
    log_ref: &Mutex<BTreeMap<usize, Out>>,
    process_log: &impl AsyncFn(Vec<Out>),
    previous_index: usize,
) -> Result<usize, usize> {
    let mut map_lock = log_ref.lock().await;
    let map = map_lock.deref_mut();
    let mut result = Vec::new();
    let mut examine_index = previous_index;
    // Rework this to peek at the first element
    while let Some((first_key, _)) = map.first_key_value() {
        let first_key = *first_key;
        if first_key == examine_index || first_key == examine_index - 1 {
            let (index, val) = map.pop_first().unwrap();
            examine_index = index;
            result.push(val);
        } else {
            break;
        }
    }
    if result.is_empty() {
        Err(previous_index)
    } else {
        process_log(result).await;
        Ok(examine_index)
    }
}

async fn async_process_table<In, Out>(
    tasks: Vec<In>,
    process_individual: &(impl AsyncFn(In) -> Out + Sync),
    log_inorder_chunk_out: &(impl AsyncFn(Vec<Out>) -> () + Sync),
    simultaneous_tasks: usize,
    shipout_every_duration: Duration,
) -> ()
where
    In: Send,
    Out: Send,
{
    let total_tasks = tasks.len();
    if total_tasks == 0 {
        return;
    }
    let tasks_deque: VecDeque<_> = tasks.into_iter().enumerate().collect();
    let tasks_mutex = Mutex::new(tasks_deque);
    let mutex_log: Mutex<BTreeMap<usize, Out>> = Mutex::new(BTreeMap::new());

    let workers = async {
        let mut worker_handles = Vec::new();
        for _ in 0..simultaneous_tasks {
            let handle = async {
                loop {
                    let task = {
                        let mut guard = tasks_mutex.lock().await;
                        guard.pop_front()
                    };
                    if let Some(t) = task {
                        wrap_individual_process(t, &mutex_log, process_individual).await;
                    } else {
                        break;
                    }
                }
            };
            worker_handles.push(handle);
        }
        futures::future::join_all(worker_handles).await;
    };

    let logger = async {
        let mut next_expected_index = 0;
        while next_expected_index < total_tasks {
            tokio::time::sleep(shipout_every_duration).await;
            match log_out_from_log(&mutex_log, log_inorder_chunk_out, next_expected_index).await {
                Ok(last_index_in_chunk) => {
                    next_expected_index = last_index_in_chunk + 1;
                }
                Err(_) => {
                    // No new items logged, continue waiting.
                }
            }
        }
    };

    tokio::join!(workers, logger);
}
