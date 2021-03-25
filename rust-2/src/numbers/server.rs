use crate::numbers::printer::Printer;
use crate::numbers::writer::Writer;
use std::collections::HashSet;
use std::io::{BufRead, BufReader};
use std::net::TcpListener;
use std::sync::atomic::{AtomicU32, AtomicU64, Ordering};
use std::sync::mpsc::{channel, sync_channel};
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

        let unique_counter = Arc::new(AtomicU32::new(0));
        let duplicate_counter = Arc::new(AtomicU64::new(0));
        let store = Arc::new(Mutex::new(HashSet::<u32>::new()));

        let (print_sender, print_receiver) = sync_channel::<String>(1);
        let p = Printer::new(
            Duration::from_secs(10),
            Arc::clone(&unique_counter),
            Arc::clone(&duplicate_counter),
        );
        p.start_print_timer(print_sender);
        thread::spawn(move || loop {
            let msg = print_receiver.recv().unwrap();
            println!("{}", msg);
        });

        let (tx, rx) = channel::<String>();

        let rx = Arc::new(Mutex::new(rx));
        let file_path = Arc::new(String::from("numbers.log"));
        let writer = Writer::new(rx, file_path);
        writer.start_log_file_output();

        // listener thread
        let listener = TcpListener::bind(addr).unwrap();
        for stream in listener.incoming() {
            let stream = stream.unwrap();
            let store = Arc::clone(&store);
            let unique_counter = Arc::clone(&unique_counter);
            let duplicate_counter = Arc::clone(&duplicate_counter);
            let tx = tx.clone();
            thread::spawn(move || {
                let mut buf = BufReader::new(stream);
                loop {
                    let mut input = String::new();
                    buf.read_line(&mut input).unwrap();
                    if input.is_empty() {
                        break;
                    }
                    let num: u32 = input.trim().parse().unwrap();
                    debug!("Got number: {}", &num);
                    if store.lock().unwrap().insert(num) {
                        unique_counter.fetch_add(1, Ordering::SeqCst);
                        tx.send(String::from(input.trim()));
                    } else {
                        duplicate_counter.fetch_add(1, Ordering::SeqCst);
                    }
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
