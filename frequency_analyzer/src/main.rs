use std::collections::HashMap;
use std::io::Read;

fn main() -> Result<(), std::io::Error> {
    let mut buffer = String::new();
    std::io::stdin().read_to_string(&mut buffer)?;

    let alphanumeric : String = buffer.chars().filter(|c| c.is_alphanumeric()).collect();

    let count_map = alphanumeric.chars().fold(HashMap::new(), |mut map, c| {
        *map.entry(c).or_insert(0) += 1;
        map
    });

    let mut count_vec: Vec<_> = count_map.iter().collect();

    count_vec.sort_by(|a, b| b.1.cmp(a.1));

    for (character,count) in count_vec.iter() {
        let percentage : f64 = (*count * 100 / alphanumeric.len()) as f64;
        println!("'{}': {} - {}%", character, *count, percentage);
    }

    Ok(())
}
