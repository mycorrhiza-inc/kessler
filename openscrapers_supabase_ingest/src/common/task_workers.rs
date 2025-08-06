use std::{
    any::{Any, TypeId},
    cmp::Ordering,
    collections::{BinaryHeap, HashMap},
    convert::Infallible,
    sync::LazyLock,
    time::Instant,
};

use async_trait::async_trait;
use rand::random;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use tokio::sync::{Mutex, RwLock, Semaphore};

struct PriorityTaskObject {
    priority: i32,
    timestamp: Instant,
    task_id: u64,
    task_object: Box<dyn ExecuteUserTask>,
}
impl PriorityTaskObject {
    pub fn new(obj: Box<dyn ExecuteUserTask>, priority: i32) -> Self {
        PriorityTaskObject {
            priority,
            task_object: obj,
            timestamp: Instant::now(),
            task_id: random(),
        }
    }
    pub fn get_task_type(&self) -> TypeId {
        (*self.task_object).type_id()
    }
    pub fn get_task_type_string(&self) -> String {
        typeid_debug(self.get_task_type())
    }
}
pub fn typeid_debug(t: TypeId) -> String {
    format!("{t:?}")
}

impl PartialEq for PriorityTaskObject {
    fn eq(&self, other: &Self) -> bool {
        self.priority == other.priority
            && self.timestamp == other.timestamp
            && self.task_id == other.task_id
    }
}
impl Eq for PriorityTaskObject {}

/* -------------------------- Ordering ------------------------------------- */

impl PartialOrd for PriorityTaskObject {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        // The total ordering is defined in `Ord::cmp`, so we just forward to it.
        Some(self.cmp(other))
    }
}

impl Ord for PriorityTaskObject {
    fn cmp(&self, other: &Self) -> Ordering {
        // 1️⃣  Higher `priority` wins → normal `i32::cmp`.
        // 2️⃣  Earlier (`lower`) timestamp wins → reverse the `Instant` ordering.
        // 3️⃣  Larger `task_id` wins → normal `u64::cmp` (deterministic tie‑breaker).
        self.priority
            .cmp(&other.priority)
            .then_with(|| other.timestamp.cmp(&self.timestamp)) // reverse
            .then_with(|| self.task_id.cmp(&other.task_id))
    }
}
// Could you derive Ord on this object, ranking higher priorities higher, and on ties rank the
// lower timestamp higher.

static TASK_PRIORITY_QUEUE: Mutex<BinaryHeap<PriorityTaskObject>> =
    Mutex::const_new(BinaryHeap::new());

static TASK_STATUS_DATA: LazyLock<RwLock<HashMap<u64, TaskStatus>>> =
    LazyLock::new(|| RwLock::new(HashMap::new()));

async fn add_task_to_queue(obj: Box<dyn ExecuteUserTask>, priority: i32) {
    let task_object = PriorityTaskObject::new(obj, priority);
    let task_id = task_object.task_id;
    let task_type_str = task_object.get_task_type_string();
    let task_status = TaskStatus {
        task_id,
        task_type_str,
        completed: false,
        state: TaskState::Waiting,
        sucess_value: None,
        error_value: None,
    };
    let mut task_status_writelock = (*TASK_STATUS_DATA).write().await;
    task_status_writelock.insert(task_id, task_status);
    drop(task_status_writelock);

    let mut queue_guard = TASK_PRIORITY_QUEUE.lock().await;
    queue_guard.push(task_object);
    drop(queue_guard);
}

async fn pop_task_from_queue() -> Option<PriorityTaskObject> {
    let mut queue_guard = TASK_PRIORITY_QUEUE.lock().await;
    let option = queue_guard.pop();
    drop(queue_guard);
    match option {
        None => None,
        Some(val) => {
            let mut task_status_writelock = (*TASK_STATUS_DATA).write().await;
            let optref = task_status_writelock.get_mut(&val.task_id);
            if let Some(mutref) = optref {
                mutref.state = TaskState::Processing
            };
            Some(val)
        }
    }
}

const SIMULTANEOUS_USER_TASKS: usize = 20;
static USER_TASK_SEMAPHORE: Semaphore = Semaphore::const_new(SIMULTANEOUS_USER_TASKS);

async fn start_workers() -> Infallible {
    loop {
        let permit = USER_TASK_SEMAPHORE.acquire().await.unwrap();
        match pop_task_from_queue().await {
            None => drop(permit),
            Some(res) => {
                tokio::spawn(async move {
                    let task_id = res.task_id;
                    let obj = res.task_object;
                    let task_completed_status = obj.execute_task(task_id).await;

                    let mut task_status_writelock = (*TASK_STATUS_DATA).write().await;
                    task_status_writelock.insert(task_id, task_completed_status);
                    drop(task_status_writelock);
                    drop(permit);
                });
            }
        }
    }
}

#[derive(Clone, Copy, Serialize, Deserialize, JsonSchema, Debug)]
enum TaskState {
    Waiting,
    Processing,
    Successful,
    Errored,
}

#[derive(Clone, Serialize, Deserialize, JsonSchema, Debug)]
pub struct TaskStatus {
    task_id: u64,
    completed: bool,
    state: TaskState,
    task_type_str: String,
    sucess_value: Option<serde_json::Value>,
    error_value: Option<serde_json::Value>,
}

#[async_trait]
pub trait ExecuteUserTask: Send {
    async fn execute_task(self: Box<Self>, task_id: u64) -> TaskStatus;
}
