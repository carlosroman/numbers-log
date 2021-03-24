use std::collections::HashSet;
use std::fs::File;
use std::io::{BufRead, BufReader, LineWriter, Write};
use std::net::TcpListener;
use std::path::Path;
use std::sync::atomic::{AtomicU32, AtomicU64, Ordering};
use std::sync::mpsc::{channel, Receiver};
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::Duration;

pub struct Server {
    port: u16,
    host: String,
}

fn start_log_file_output(rx: Receiver<String>) {
    // Create a path to the desired file
    let path = Path::new("numbers.log");
    let display = path.display();
    thread::spawn(move || {
        let file = match File::create(&path) {
            Err(why) => panic!("could not create {}: {}", display, why),
            Ok(file) => file,
        };
        let mut file = LineWriter::new(file);

        loop {
            match rx.recv() {
                Ok(val) => {
                    file.write_all(val.as_bytes());
                    file.write_all(b"\n");
                }
                Err(e) => {
                    error!("Got following error: {:?}", e);
                    return;
                }
            }
        }
    });
}

fn start_print_timer(unique_counter: &Arc<AtomicU32>, duplicate_counter: &Arc<AtomicU64>) {
    let unique_counter = Arc::clone(&unique_counter);
    let duplicate_counter = Arc::clone(&duplicate_counter);

    // timer thread
    thread::spawn(move || {
        let interval = Duration::from_secs(10);
        let mut last_unique_count = 0;
        let mut last_duplicate_count = 0;
        loop {
            thread::sleep(interval);
            let unique_count = u64::from(unique_counter.load(Ordering::SeqCst));
            let duplicate_count = duplicate_counter.load(Ordering::SeqCst);
            println!(
                "Received {} unique numbers, {} duplicates. Unique total: {}",
                unique_count - last_unique_count,
                duplicate_count - last_duplicate_count,
                unique_count
            );
            last_unique_count = unique_count;
            last_duplicate_count = duplicate_count;
        }
    });
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

        start_print_timer(&unique_counter, &duplicate_counter);

        let (tx, rx) = channel::<String>();
        start_log_file_output(rx);

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
