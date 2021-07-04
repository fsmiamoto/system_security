use std::collections::HashMap;
use std::io::Read;

const ASCII_LOWER_LIMIT: u8 = 31;
const ASCII_UPPER_LIMIT: u8 = 127;

fn main() -> Result<(), std::io::Error> {
    let mut buffer = Vec::new();
    std::io::stdin().read_to_end(&mut buffer)?;

    let sanitized: Vec<char> = buffer
        .into_iter()
        .filter(|c| (*c > ASCII_LOWER_LIMIT && *c < ASCII_UPPER_LIMIT))
        .map(|c| (c as char))
        .collect();

    let alphanumeric: String = sanitized.iter().filter(|c| c.is_alphanumeric()).collect();

    let count_map = alphanumeric.chars().fold(HashMap::new(), |mut map, c| {
        *map.entry(c).or_insert(0) += 1;
        map
    });

    let mut count_vec: Vec<_> = count_map.iter().collect();

    count_vec.sort_by(|a, b| b.1.cmp(a.1));

    for (character, count) in count_vec.iter() {
        let percentage: f64 = (*count * 100 / alphanumeric.len()) as f64;
        println!("'{}': {} - {}%", character, *count, percentage);
    }

    Ok(())
}
