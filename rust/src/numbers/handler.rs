use crate::numbers::store::Store;
use std::cell::RefCell;
use std::io::{BufRead, BufReader};
use std::sync::Arc;

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
}

pub struct StoreHandler {
    store: Arc<RefCell<Box<dyn Store>>>,
}

impl StoreHandler {
    pub fn new(store: Arc<RefCell<Box<dyn Store>>>) -> Arc<StoreHandler> {
        let h = StoreHandler { store };
        Arc::new(h)
    }
}

impl Handler for StoreHandler {
    fn handle_client(&self, stream: impl std::io::Read) {
        let mut buf = BufReader::new(stream);
        let mut input = String::new();
        buf.read_line(&mut input).unwrap();
        let num: u32 = input.trim().parse().unwrap();
        let store = Arc::clone(&self.store);
        store.borrow_mut().save(num);
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
}
