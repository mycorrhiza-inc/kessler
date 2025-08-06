use std::{
    collections::{BinaryHeap, HashMap},
    convert::Infallible,
    ops::DerefMut,
    sync::LazyLock,
    time::Instant,
};

use async_trait::async_trait;
use rand::random;
use tokio::sync::{Mutex, RwLock, Semaphore};

struct PriorityTaskObject {
    priority: i32,
    timestamp: Instant,
    task_id: u64,
    task_object: Box<dyn ExecuteUserTask>,
}

// Could you derive Ord on this object, ranking higher priorities higher, and on ties rank the
// lower timestamp higher.

static TASK_PRIORITY_QUEUE: Mutex<BinaryHeap<PriorityTaskObject>> =
    Mutex::const_new(BinaryHeap::new());

async fn add_task_to_queue(obj: Box<dyn ExecuteUserTask>, priority: i32) {
    let task_id = random();
    let task_object = PriorityTaskObject {
        priority,
        task_object: obj,
        task_id,
        timestamp: Instant::now(),
    };
    let queue_guard = TASK_PRIORITY_QUEUE.lock().await;
    let queue_ref = queue_guard.deref_mut();
    queue_ref.push(task_object);
}

async fn pop_task_from_queue() -> Option<PriorityTaskObject> {
    let queue_guard = TASK_PRIORITY_QUEUE.lock().await;
    let queue_ref = queue_guard.deref_mut();
    let option = queue_ref.pop();
    option
}

const SIMULTANEOUS_USER_TASKS: usize = 20;
static USER_TASK_SEMAPHORE: Semaphore = Semaphore::const_new(SIMULTANEOUS_USER_TASKS);

async fn start_workers() -> Infallible {
    loop {
        let permit = USER_TASK_SEMAPHORE.acquire().await.unwrap();
        match pop_task_from_queue() {
            None => drop(permit),
            Some(res) => tokio::spawn(async move {
                let obj = res.obj;
            }),
        }
    }
}

static TASK_STATUS_DATA: LazyLock<RwLock<HashMap<u64, TaskStatus>>> =
    LazyLock::new(|| RwLock::new(HashMap::new()));

#[derive(Clone, Copy)]
enum TaskState {
    Waiting,
    Processing,
    Successful,
    Errored,
}

struct TaskStatus {
    completed: bool,
    error: Option<String>,
}

#[async_trait]
trait ExecuteUserTask: Send {
    async fn execute_task(self: Box<Self>) -> TaskStatus;
}
