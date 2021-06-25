use std::io::Read;
use structopt::StructOpt;
mod caesar;

#[derive(Debug, StructOpt)]
#[structopt(name = "caesar", about = "A simple Caesar Cipher")]
struct Cli {
    #[structopt(short, long)]
    cipher: bool,
    #[structopt(short, long)]
    decipher: bool,
    #[structopt(short = "k", long = "key")]
    key: u8,
}

fn main() -> Result<(), &'static str> {
    let args = Cli::from_args();

    if args.cipher && args.decipher {
        return Err("cipher and decipher flags are mutually exclusive");
    }

    if !args.cipher && !args.decipher {
        return Err("at least one flag is required");
    }

    let mut buffer = String::new();
    let result = std::io::stdin().read_to_string(&mut buffer);
    match result {
        Ok(_) => {}
        Err(_) => return Err("bad line"),
    };

    let result = if args.cipher {
        caesar::cipher(buffer, args.key)
    } else {
        caesar::decipher(buffer, args.key)
    };
    println!("{}", result);
    Ok(())
}
