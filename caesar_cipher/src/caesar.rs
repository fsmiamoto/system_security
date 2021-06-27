// TODO: Remove magic numbers

const ALPHABET: [char; 62] = [
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
    'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
    'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4',
    '5', '6', '7', '8', '9',
];

pub fn cipher(content: String, key: u8) -> String {
    content
        .chars()
        .map(|c| remove_diacritics(c))
        .map(|c| cipher_char(c, key))
        .collect()
}

pub fn decipher(content: String, key: u8) -> String {
    content
        .chars()
        .map(|c| remove_diacritics(c))
        .map(|c| decipher_char(c, key))
        .collect()
}

fn cipher_char(c: char, key: u8) -> char {
    if !c.is_alphanumeric() {
        return c;
    }
    let index = get_index(c);
    let other_index = if index + key < 62 {
        index + key
    } else {
        index + key - 62
    };
    ALPHABET[other_index as usize]
}

fn decipher_char(c: char, key: u8) -> char {
    if !c.is_alphanumeric() {
        return c;
    }
    let index = get_index(c);
    let other_index = if index >= key {
        index - key
    } else {
        index + 62 - key
    };
    ALPHABET[other_index as usize]
}

fn get_index(c: char) -> u8 {
    if c.is_uppercase() {
        return (c as u32 - 'A' as u32) as u8;
    }
    if c.is_lowercase() {
        return (c as u32 - 'a' as u32) as u8 + 26;
    }
    if c.is_numeric() {
        return (c as u32 - '0' as u32) as u8 + 52;
    }
    return 0;
}

fn remove_diacritics(c: char) -> char {
    match c {
        'á' | 'à' | 'ã' => 'a',
        'é' | 'ê' | 'è' => 'e',
        'ó' | 'ô' => 'o',
        'í' => 'i',
        'ú' => 'u',
        'ç' => 'c',

        'Á' | 'À' | 'Ã' => 'A',
        'É' | 'Ê' => 'E',
        'Ó' | 'Ô' => 'O',
        'Í' => 'I',
        'Ú' => 'U',
        'Ç' => 'C',

        _ => c,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cipher() -> Result<(), String> {
        let key = 23;
        let message: String = "THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG".to_string();
        assert_eq!(
            cipher(message, key),
            "qeb nrfZh Yoltk clu grjmp lsbo qeb iXwv ald"
        );
        Ok(())
    }

    #[test]
    fn test_decipher() -> Result<(), String> {
        let key = 23;
        let message: String = "qeb nrfZh Yoltk clu grjmp lsbo qeb iXwv ald".to_string();
        assert_eq!(
            decipher(message, key),
            "THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG"
        );
        Ok(())
    }

    #[test]
    fn test_remove_diacritics() -> Result<(), String> {
        assert_eq!(remove_diacritics('é'), 'e');
        assert_eq!(remove_diacritics('ã'), 'a');
        Ok(())
    }

    #[test]
    fn test_decipher_char_edges() -> Result<(), String> {
        let key = 2;

        assert_eq!(decipher_char('a', key), 'Y');
        assert_eq!(decipher_char('0', key), 'y');
        assert_eq!(decipher_char('A', key), '8');

        Ok(())
    }

    #[test]
    fn test_cipher_char_edges() -> Result<(), String> {
        let key = 1;

        assert_eq!(decipher_char('a', key), 'Z');
        assert_eq!(decipher_char('0', key), 'z');
        assert_eq!(decipher_char('A', key), '9');

        Ok(())
    }

    #[test]
    fn test_decipher_char_wrap() -> Result<(), String> {
        let key = 5;
        assert_eq!(decipher_char('A', key), '5');

        Ok(())
    }

    #[test]
    fn test_cipher_char_wrap() -> Result<(), String> {
        let key = 5;

        assert_eq!(cipher_char('9', key), 'E');

        Ok(())
    }
}
