use structopt::StructOpt;

#[derive(Debug, StructOpt)]
#[structopt(name="caesar", about = "A simple Caesar Cipher")]
struct Args {
    #[structopt(short,long)]
    cipher: bool,
    #[structopt(short,long)]
    decipher: bool,
    #[structopt(short = "k", long ="key")]
    key: i8
}

fn main() {
    let args = Args::from_args();
    println!("{:?}", args);
}

