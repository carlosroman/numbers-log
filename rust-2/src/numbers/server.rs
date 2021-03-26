use crate::numbers::printer::Printer;
use crate::numbers::store::BTreeSetStore;
use crate::numbers::store::Store;
use crate::numbers::writer::Writer;
use std::io::{BufRead, BufReader};
use std::net::TcpListener;
use std::sync::mpsc::sync_channel;
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::Duration;

pub struct Server {
    port: u16,
    host: String,
}

impl Server {
    pub fn new(host: String, port: u16) -> Server {
        Server { host, port }
    }

    pub fn start(&self) {
        let addr = format!("{}:{}", self.host, self.port);
        info!("Starting server at: {}", addr);

        // Setup log writer
        let (writer_sender, writer_receiver) = sync_channel::<String>(5 * 1000);
        let writer_receiver = Arc::new(Mutex::new(writer_receiver));
        let file_path = Arc::new(String::from("numbers.log"));
        let writer = Writer::new(writer_receiver, file_path);
        writer.start_log_file_output();

        // Setup store
        let store = BTreeSetStore::new();
        let sender = store.start_processing(5 * 1000, writer_sender);
        let duplicate_counter = store.duplicate_counter();
        let unique_counter = store.unique_counter();

        // Setup stats printer
        let (print_sender, print_receiver) = sync_channel::<String>(1);
        let p = Printer::new(
            Duration::from_secs(10),
            Arc::clone(&unique_counter),
            Arc::clone(&duplicate_counter),
        );
        p.start_stats_timer(print_sender);
        thread::spawn(move || loop {
            let msg = print_receiver.recv().unwrap();
            println!("{}", msg);
        });

        // listener thread
        let listener = TcpListener::bind(addr).unwrap();
        for stream in listener.incoming() {
            let stream = stream.unwrap();
            let num_sender = sender.clone();
            thread::spawn(move || {
                let mut buf = BufReader::new(stream);
                loop {
                    let mut input = String::new();
                    buf.read_line(&mut input).unwrap();
                    if input.is_empty() {
                        break;
                    }
                    let num: u32 = input.trim().parse().unwrap();
                    num_sender.send(num).unwrap();
                }
            });
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::net::TcpListener;
    use std::thread;

    fn init() {
        let _ = env_logger::builder().is_test(true).try_init();
    }

    #[test]
    fn server_passes_bytes_to_processor() {
        init();
        let s = Server::new(String::from("0.0.0.0"), get_free_port().unwrap());

        thread::spawn(move || {
            s.start();
        });
    }

    fn get_free_port() -> Result<u16, String> {
        for port in 1025..65535 {
            match TcpListener::bind(("0.0.0.0", port)) {
                Ok(_l) => {
                    return Ok(port);
                }
                _ => {}
            }
        }
        Err(String::from("no port found"))
    }
}
