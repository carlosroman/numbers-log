use std::fs::File;
use std::io::{LineWriter, Write};
use std::path::Path;
use std::sync::mpsc::Receiver;
use std::sync::{Arc, Mutex};
use std::thread;

pub struct Writer {
    file_path: Arc<String>,
    rx: Arc<Mutex<Receiver<String>>>,
}

impl Writer {
    pub fn new(rx: Arc<Mutex<Receiver<String>>>, file_path: Arc<String>) -> Writer {
        Writer { rx, file_path }
    }

    pub fn start_log_file_output(&self) {
        let file_path = Arc::clone(&self.file_path);
        let rx = Arc::clone(&self.rx);
        thread::spawn(move || {
            // let file_path = Arc::clone(&file_path);
            let path = Path::new(file_path.as_str());
            let display = path.display();
            let file = match File::create(&path) {
                Err(why) => panic!("could not create {}: {}", display, why),
                Ok(file) => file,
            };
            let mut file = LineWriter::new(file);

            loop {
                match rx.lock().unwrap().recv() {
                    Ok(val) => {
                        file.write_all(val.as_bytes()).unwrap();
                        file.write_all(b"\n").unwrap();
                    }
                    Err(e) => {
                        error!("Got following error: {:?}", e);
                        return;
                    }
                }
            }
        });
    }

    pub fn stop(&self) {}
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::env::temp_dir;
    use std::io::Read;
    use std::sync::mpsc::channel;
    use std::sync::Mutex;
    use std::{thread, time};

    #[test]
    fn server_passes_bytes_to_processor() {
        // Given a temp dir
        let mut dir = temp_dir();
        dir.push("temp");

        let (tx, rx) = channel::<String>();
        let rx = Arc::new(Mutex::new(rx));
        let path = dir.to_str().unwrap();
        let file_path = Arc::new(String::from(path));

        // And a writer is setup to write to temp file
        let w = Writer::new(rx, file_path);
        w.start_log_file_output();

        // When I send a message to write
        assert!(tx.send(String::from("bob")).is_ok());

        // Then I expect to see that message in the temp file
        thread::sleep(time::Duration::from_millis(10));
        let mut result = File::open(&path).unwrap();
        let mut buffer = String::new();
        result.read_to_string(&mut buffer).unwrap();
        assert_eq!("bob\n", String::from(buffer));
    }
}
