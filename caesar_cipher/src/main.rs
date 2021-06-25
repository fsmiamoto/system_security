use std::char;
use std::io::Read;
use structopt::StructOpt;

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

    let operator: fn(char, u8) -> char = if args.cipher { cipher } else { decipher };

    let result: String = buffer.chars().map(|c| operator(c, args.key)).collect();
    println!("{}", result);
    Ok(())
}

// TODO: Check for corner cases
fn cipher(c: char, key: u8) -> char {
    if !c.is_alphanumeric() {
        return c;
    }
    return std::char::from_u32(c as u32 + key as u32).unwrap_or(c);
}

fn decipher(c: char, key: u8) -> char {
    if !c.is_alphanumeric() {
        return c;
    }
    return std::char::from_u32(c as u32 - key as u32).unwrap_or(c);
}
