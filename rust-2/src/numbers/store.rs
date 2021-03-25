use std::collections::HashSet;
use std::sync::atomic::{AtomicU32, AtomicU64, Ordering};
use std::sync::mpsc::{sync_channel, SyncSender};
use std::sync::Arc;
use std::thread;

const MAX_NUMBER: usize = 1000000000;

pub struct Store {
    unique_counter: Arc<AtomicU32>,
    duplicate_counter: Arc<AtomicU64>,
}

impl Store {
    pub fn new() -> Store {
        let unique_counter = Arc::new(AtomicU32::new(0));
        let duplicate_counter = Arc::new(AtomicU64::new(0));
        Store {
            unique_counter,
            duplicate_counter,
        }
    }

    pub fn start_processing(
        &self,
        buffer_size: usize,
        resp_tx: SyncSender<String>,
    ) -> SyncSender<u32> {
        let (tx, rx) = sync_channel::<u32>(buffer_size);
        let unique_counter = Arc::clone(&self.unique_counter);
        let duplicate_counter = Arc::clone(&self.duplicate_counter);
        thread::spawn(move || {
            let mut store = HashSet::<u32>::with_capacity(MAX_NUMBER);
            loop {
                let res = rx.recv().unwrap();
                if store.insert(res) {
                    unique_counter.fetch_add(1, Ordering::SeqCst);
                    resp_tx.send(format!("{:09}", res)).unwrap();
                } else {
                    duplicate_counter.fetch_add(1, Ordering::SeqCst);
                }
            }
        });
        return tx;
    }

    pub fn duplicate_counter(&self) -> Arc<AtomicU64> {
        self.duplicate_counter.clone()
    }

    pub fn unique_counter(&self) -> Arc<AtomicU32> {
        self.unique_counter.clone()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::atomic::Ordering;
    use std::time::Duration;

    #[test]
    fn store_saves_values() {
        let s = Store::new();

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
