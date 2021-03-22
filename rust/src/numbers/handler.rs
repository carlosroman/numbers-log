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
