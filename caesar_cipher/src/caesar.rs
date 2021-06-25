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
    std::char::from_u32(c as u32 + key as u32).unwrap_or(c)
}

fn decipher_char(c: char, key: u8) -> char {
    if !c.is_alphanumeric() {
        return c;
    }
    std::char::from_u32(c as u32 - key as u32).unwrap_or(c)
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
    fn test_remove_diacritics() -> Result<(), String> {
        assert_eq!(remove_diacritics('é'), 'e');
        assert_eq!(remove_diacritics('ã'), 'a');
        Ok(())
    }

    #[test]
    fn test_decipher_char() -> Result<(), String> {
        let key = 1;

        // Basic
        assert_eq!(decipher_char('b', key), 'a');
        assert_eq!(decipher_char('9', key), '8');

        // Edge cases
        assert_eq!(decipher_char('a', key), 'Z');
        assert_eq!(decipher_char('0', key), 'z');
        assert_eq!(decipher_char('A', key), '9');

        Ok(())
    }
}
