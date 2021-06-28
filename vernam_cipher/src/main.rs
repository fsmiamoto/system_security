use std::io::Read;
use std::path::PathBuf;
use structopt::StructOpt;

#[derive(StructOpt)]
#[structopt(name = "vernam", about = "A simple Vernam Cipher")]
struct Cli {
    #[structopt(short ="k", long = "key",parse(from_os_str))]
    key: PathBuf,
}

fn main() -> Result<(), &'static str> {
    let args = Cli::from_args();

    let mut buffer = String::new();
    let result = std::io::stdin().read_to_string(&mut buffer);
    match result {
        Ok(_) => {}
        Err(_) => return Err("bad line"),
    };

    let content: String;
    match std::fs::read_to_string(args.key) {
        Ok(c) => { content = c}
        Err(_) => return Err("Error while reading key file")
    }

    if content.len() != buffer.len() {
        return Err("the key length must be equal to the content length");
    }

    print!("{:?}", content);
    Ok(())
}
