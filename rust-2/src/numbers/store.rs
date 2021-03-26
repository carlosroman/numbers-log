use std::collections::{BTreeSet, HashSet};
use std::sync::atomic::{AtomicU32, AtomicU64, Ordering};
use std::sync::mpsc::{sync_channel, SyncSender};
use std::sync::Arc;
use std::thread;

const MAX_NUMBER: usize = 1000000000;

pub trait Store {
    fn start_processing(&self, buffer_size: usize, resp_tx: SyncSender<String>) -> SyncSender<u32> {
        let (tx, rx) = sync_channel::<u32>(buffer_size);
        let unique_counter = Arc::clone(&self.unique_counter());
        let duplicate_counter = Arc::clone(&self.duplicate_counter());
        let mut store = self.get_store();
        thread::spawn(move || loop {
            let res = rx.recv().unwrap();
            if store.insert(res) {
                unique_counter.fetch_add(1, Ordering::SeqCst);
                resp_tx.send(format!("{:09}", res)).unwrap();
            } else {
                duplicate_counter.fetch_add(1, Ordering::SeqCst);
            }
        });
        tx
    }
    fn duplicate_counter(&self) -> Arc<AtomicU64>;
    fn unique_counter(&self) -> Arc<AtomicU32>;
    fn get_store(&self) -> Box<dyn HasInsert + Send + 'static>;
}

pub struct HashSetStore {
    unique_counter: Arc<AtomicU32>,
    duplicate_counter: Arc<AtomicU64>,
}

impl HashSetStore {
    pub fn new() -> Box<HashSetStore> {
        let unique_counter = Arc::new(AtomicU32::new(0));
        let duplicate_counter = Arc::new(AtomicU64::new(0));
        Box::new(HashSetStore {
            unique_counter,
            duplicate_counter,
        })
    }
}

pub struct BTreeSetStore {
    unique_counter: Arc<AtomicU32>,
    duplicate_counter: Arc<AtomicU64>,
}

impl BTreeSetStore {
    pub fn new() -> Box<BTreeSetStore> {
        let unique_counter = Arc::new(AtomicU32::new(0));
        let duplicate_counter = Arc::new(AtomicU64::new(0));
        Box::new(BTreeSetStore {
            unique_counter,
            duplicate_counter,
        })
    }
}

pub trait HasInsert {
    fn insert(&mut self, value: u32) -> bool;
}

impl HasInsert for HashSet<u32> {
    fn insert(&mut self, value: u32) -> bool {
        self.insert(value)
    }
}

impl HasInsert for BTreeSet<u32> {
    fn insert(&mut self, value: u32) -> bool {
        self.insert(value)
    }
}

impl Store for HashSetStore {
    fn duplicate_counter(&self) -> Arc<AtomicU64> {
        self.duplicate_counter.clone()
    }

    fn unique_counter(&self) -> Arc<AtomicU32> {
        self.unique_counter.clone()
    }

    fn get_store(&self) -> Box<dyn HasInsert + Send + 'static> {
        let store = HashSet::<u32>::with_capacity(MAX_NUMBER);
        Box::new(store)
    }
}

impl Store for BTreeSetStore {
    fn duplicate_counter(&self) -> Arc<AtomicU64> {
        self.duplicate_counter.clone()
    }

    fn unique_counter(&self) -> Arc<AtomicU32> {
        self.unique_counter.clone()
    }

    fn get_store(&self) -> Box<dyn HasInsert + Send + 'static> {
        let store = BTreeSet::<u32>::new();
        Box::new(store)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::atomic::Ordering;
    use std::time::Duration;

    #[test]
    fn hash_set_store_saves_values() {
        let s = HashSetStore::new();
        store_test(s);
    }

    #[test]
    fn btree_set_store_saves_values() {
        let s = BTreeSetStore::new();
        store_test(s);
    }

    fn store_test(s: Box<dyn Store>) {
        let (resp_tx, resp_rx) = sync_channel::<String>(0);
        let sender = s.start_processing(0, resp_tx);
        assert!(sender.send(9).is_ok());
        assert_eq!("000000009", resp_rx.recv().unwrap());

        assert_eq!(1, s.unique_counter().load(Ordering::Acquire));
        assert_eq!(0, s.duplicate_counter().load(Ordering::Acquire));

        assert!(sender.send(9).is_ok());
        assert!(resp_rx.recv_timeout(Duration::from_millis(100)).is_err());

        assert_eq!(1, s.unique_counter().load(Ordering::Acquire));
        assert_eq!(1, s.duplicate_counter().load(Ordering::Acquire));
    }
}
