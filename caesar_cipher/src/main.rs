use structopt::StructOpt;

#[derive(Debug, StructOpt)]
#[structopt(name="caesar", about = "A simple Caesar Cipher")]
struct Cli {
    #[structopt(short,long)]
    cipher: bool,
    #[structopt(short,long)]
    decipher: bool,
    #[structopt(short = "k", long ="key")]
    key: i8
}

fn main() -> Result<(), &'static str>{
    let args = Cli::from_args();

    if args.cipher && args.decipher {
        return Err("cipher and decipher flags are mutually exclusive")
    }

    println!("{:?}", args);
    Ok(())
}
