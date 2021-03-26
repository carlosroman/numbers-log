use std::sync::atomic::{AtomicU32, AtomicU64, Ordering};
use std::sync::mpsc::SyncSender;
use std::sync::Arc;
use std::thread;
use std::time::Duration;

pub struct Printer {
    interval: Duration,
    unique_counter: Arc<AtomicU32>,
    duplicate_counter: Arc<AtomicU64>,
}

impl Printer {
    pub fn new(
        interval: Duration,
        unique_counter: Arc<AtomicU32>,
        duplicate_counter: Arc<AtomicU64>,
    ) -> Printer {
        Printer {
            interval,
            unique_counter,
            duplicate_counter,
        }
    }

    pub fn start_stats_timer(&self, tx: SyncSender<String>) {
        let unique_counter = Arc::clone(&self.unique_counter);
        let duplicate_counter = Arc::clone(&self.duplicate_counter);
        let interval = self.interval;
        thread::spawn(move || {
            let mut last_unique_count = 0;
            let mut last_duplicate_count = 0;
            loop {
                thread::sleep(interval);
                let unique_count = u64::from(unique_counter.load(Ordering::SeqCst));
                let duplicate_count = duplicate_counter.load(Ordering::SeqCst);
                let res = format!(
                    "Received {} unique numbers, {} duplicates. Unique total: {}",
                    unique_count - last_unique_count,
                    duplicate_count - last_duplicate_count,
                    unique_count
                );
                tx.send(res).unwrap();
                last_unique_count = unique_count;
                last_duplicate_count = duplicate_count;
            }
        });
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::mpsc::sync_channel;

    #[test]
    fn printer_outputs_data() {
        let unique_counter = Arc::new(AtomicU32::new(13));
        let duplicate_counter = Arc::new(AtomicU64::new(37));

        let (sender, receiver) = sync_channel::<String>(1);
        let p = Printer::new(Duration::from_millis(10), unique_counter, duplicate_counter);
        p.start_stats_timer(sender.clone());

        let actual = receiver.recv().unwrap();
        assert_eq!(
            "Received 13 unique numbers, 37 duplicates. Unique total: 13",
            actual
        );
    }
}
