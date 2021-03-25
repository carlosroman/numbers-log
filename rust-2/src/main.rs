#[macro_use]
extern crate log;

use crate::numbers::server::Server;

mod numbers {
    pub mod server;
    pub mod writer;
    pub mod printer;
    pub mod store;
}

fn main() {
    env_logger::init();
    info!("Starting");
    let s = Server::new(String::from("0.0.0.0"), 4000);
    s.start();
}
