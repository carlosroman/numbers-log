use crate::numbers::handler;
use crate::numbers::server::Server;

mod numbers {
    pub mod handler;
    pub mod server;
    pub mod store;
}

fn main() {
    let _s = Server::new(String::from("0.0.0.0"), 4000, handler::NoopHandler::new());
}
