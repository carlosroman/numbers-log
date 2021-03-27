use std::sync::{Arc, Mutex};

use crate::numbers::handler;
use crate::numbers::server::Server;
use crate::numbers::store;

mod numbers {
    pub mod handler;
    pub mod server;
    pub mod store;
}

fn main() {
    let store = Arc::new(Mutex::new(store::InMemoryStore::new()));
    let handler = handler::StoreHandler::new(store);
    let s = Server::new(String::from("0.0.0.0"), 4000, handler);
    s.start();
}
