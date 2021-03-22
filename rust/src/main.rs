use crate::numbers::server::Server;
use crate::numbers::handler;

mod numbers {
    pub mod server;
    pub mod handler;
}

fn main() {
    let _s = Server::new(
        String::from("0.0.0.0"),
        4000,
        handler::NoopHandler::new());
}
