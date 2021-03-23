use crate::numbers::store::Store;
use std::cell::RefCell;
use std::io::{BufRead, BufReader};
use std::sync::{Arc, Mutex};

pub struct NoopHandler {}

impl NoopHandler {
    pub fn new() -> Arc<NoopHandler> {
        let h = NoopHandler {};
        Arc::new(h)
    }
}

impl Handler for NoopHandler {
    fn handle_client(&self, _stream: impl std::io::Read) {}
}

pub trait Handler {
    fn handle_client(&self, _stream: impl std::io::Read) {}

    fn get_stats(&self) -> (u32, u64) {
        (u32::max_value(), u64::max_value())
    }
}

pub struct StoreHandler {
    store: Arc<RefCell<Box<dyn Store>>>,
    unique_counter: Arc<Mutex<u32>>,
    duplicate_counter: Arc<Mutex<u64>>,
}

impl StoreHandler {
    pub fn new(store: Arc<RefCell<Box<dyn Store>>>) -> Arc<StoreHandler> {
        let unique_counter = Arc::new(Mutex::new(0));
        let duplicate_counter = Arc::new(Mutex::new(0));
        let h = StoreHandler {
            store,
            unique_counter,
            duplicate_counter,
        };
        Arc::new(h)
    }
}

impl Handler for StoreHandler {
    fn handle_client(&self, stream: impl std::io::Read) {
        let mut buf = BufReader::new(stream);

        loop {
            let mut input = String::new();
            buf.read_line(&mut input).unwrap();
            if input.is_empty() {
                break;
            }

            let num: u32 = input.trim().parse().unwrap();
            let store = Arc::clone(&self.store);
            match store.borrow_mut().save(num) {
                Some(res) => {
                    if res {
                        let counter = Arc::clone(&self.unique_counter);
                        let mut num = counter.lock().unwrap();
                        *num += 1;
                    } else {
                        let counter = Arc::clone(&self.duplicate_counter);
                        let mut num = counter.lock().unwrap();
                        *num += 1;
                    }
                }
                _ => {
                    break;
                }
            };
        }
    }

    fn get_stats(&self) -> (u32, u64) {
        let unique = *Arc::clone(&self.unique_counter).lock().unwrap();
        let duplicate = *Arc::clone(&self.duplicate_counter).lock().unwrap();
        (unique, duplicate)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::any::Any;
    use std::collections::HashSet;

    struct StoreStub {
        store: HashSet<u32>,
    }

    impl StoreStub {
        fn new() -> StoreStub {
            StoreStub {
                store: HashSet::new(),
            }
        }
    }

    impl Store for StoreStub {
        fn save(&mut self, val: u32) -> Option<bool> {
            Some(self.store.insert(val))
        }

        fn as_any(&self) -> &dyn Any {
            self
        }
    }

    #[test]
    fn store_handler_saves_parsed_value_to_store() {
        let bs: Box<dyn Store> = Box::new(StoreStub::new());
        let s = Arc::new(RefCell::new(bs));
        let h = StoreHandler::new(s.clone());
        let text = "999999999\n".as_bytes();
        h.handle_client(text);

        match &s.borrow().as_any().downcast_ref::<StoreStub>() {
            Some(b) => {
                assert!(b.store.contains(&999999999));
            }
            None => panic!("Did not find a store stub"),
        };
    }

    #[test]
    fn store_handler_increments_unique_counter() {
        let h = StoreHandler::new(Arc::new(RefCell::new(Box::new(StoreStub::new()))));
        let text = "999999999\n".as_bytes();
        h.handle_client(text);
        let (unique, duplicate) = h.get_stats();
        assert_eq!(1, unique);
        assert_eq!(0, duplicate);
        assert_eq!(1, duplicate + u64::from(unique)); // total
    }

    #[test]
    fn store_handler_increments_duplicate_counter() {
        let h = StoreHandler::new(Arc::new(RefCell::new(Box::new(StoreStub::new()))));
        let text = "999999999\n999999999\n".as_bytes();
        h.handle_client(text);
        let (unique, duplicate) = h.get_stats();
        assert_eq!(1, unique);
        assert_eq!(1, duplicate);
        assert_eq!(2, duplicate + u64::from(unique)); // total
    }
}
