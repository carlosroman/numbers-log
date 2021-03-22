use crate::numbers::handler::Handler;
use std::net::TcpListener;
use std::sync::Arc;

pub struct Server<T: Handler> {
    handler: Arc<T>,
    host: String,
    port: u16,
}

impl<T: Handler> Server<T> {
    pub fn new(host: String, port: u16, handler: Arc<T>) -> Server<T> {
        Server {
            host,
            port,
            handler,
        }
    }

    pub fn start(&self) {
        let addr = format!("{}:{}", self.host, self.port);
        println!("{}", addr.clone());
        let listener = TcpListener::bind(addr).unwrap();
        for stream in listener.incoming() {
            let handler = Arc::clone(&self.handler);
            match stream {
                Ok(stream) => {
                    handler.handle_client(stream);
                }
                Err(e) => {
                    println!("Error: {}", e);
                }
            }
        }
    }

    // pub fn stop(&self) {}
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::numbers::handler;
    use std::io::{BufRead, BufReader, Write};
    use std::net::{TcpListener, TcpStream};
    use std::sync::mpsc::{sync_channel, Receiver, SyncSender};
    use std::sync::Arc;
    use std::thread;

    struct MockHandler {
        tx: SyncSender<String>,
    }

    impl MockHandler {
        fn new() -> (MockHandler, Receiver<String>) {
            let (tx, rx) = sync_channel::<String>(2);
            return (MockHandler { tx }, rx);
        }
    }

    impl handler::Handler for MockHandler {
        fn handle_client(&self, stream: impl std::io::Read) {
            let mut buf = BufReader::new(stream);
            let mut input = String::new();
            buf.read_line(&mut input).unwrap();
            println!("Read the following: {}", input.trim());
            self.tx.send(input).unwrap();
        }
    }

    #[test]
    fn server_passes_bytes_to_processor() {
        let (mock_handler, rx) = MockHandler::new();
        let mock_handler = Arc::new(mock_handler);
        let server_handler = Arc::clone(&mock_handler);
        let port = get_free_port().unwrap();
        let server = Server::new(String::from("localhost"), port.clone(), server_handler);

        thread::spawn(move || {
            server.start();
        });

        // thread::sleep(time::Duration::from_secs(5));
        match TcpStream::connect(format!("localhost:{}", port)) {
            Ok(mut s) => {
                let _mock_server = Arc::clone(&mock_handler);
                let res = s.write(b"Expected Text\n");
                assert!(!res.is_err());
                let actual = rx.recv().unwrap();
                assert_eq!(actual, "Expected Text\n");
            }
            Err(e) => {
                panic!("Test failed as got error: {}", e);
            }
        }
    }

    fn get_free_port() -> Result<u16, String> {
        for port in 1025..65535 {
            match TcpListener::bind(("127.0.0.1", port)) {
                Ok(_l) => {
                    return Ok(port);
                }
                _ => {}
            }
        }
        Err(String::from("no port found"))
    }
}
