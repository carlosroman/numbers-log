use std::any::Any;
use std::collections::HashSet;
use std::sync::{Arc, Mutex};

pub trait Store {
    fn save(&mut self, _val: u32) -> Option<bool> {
        None
    }
    fn as_any(&self) -> &dyn Any;
}

pub struct InMemoryStore {
    store: Arc<Mutex<HashSet<u32>>>,
}

impl InMemoryStore {
    fn new() -> Box<dyn Store> {
        let store = Arc::new(Mutex::new(HashSet::new()));
        Box::new(InMemoryStore { store })
    }
}

impl Store for InMemoryStore {
    fn save(&mut self, val: u32) -> Option<bool> {
        let store = Arc::clone(&self.store);
        let mut guard = store.lock().unwrap();
        Some(guard.insert(val))
    }

    fn as_any(&self) -> &dyn Any {
        self
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn store_returns_true_if_number_has_not_been_seen() {
        let mut s = InMemoryStore::new();
        let actual = s.save(1).unwrap();
        assert_eq!(actual, true);
    }

    #[test]
    fn store_returns_false_if_number_has_been_seen() {
        let mut s = InMemoryStore::new();
        s.save(1);
        let actual = s.save(1).unwrap();
        assert_eq!(actual, false);
    }
}
