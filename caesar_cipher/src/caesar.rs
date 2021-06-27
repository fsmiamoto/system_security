const ALPHABET: [char; 62] = [
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
    'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l',
    'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4',
    '5', '6', '7', '8', '9',
];

const ALPHABET_LEN: u8 = 62;
const UPPER_BASE_INDEX: u8 = 0;
const LOWER_BASE_INDEX: u8 = 26;
const DIGIT_BASE_INDEX: u8 = 52;

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
    let before = get_index(c);
    let after = if before + key < ALPHABET_LEN {
        before + key
    } else {
        before + key - ALPHABET_LEN
    };
    ALPHABET[after as usize]
}

fn decipher_char(c: char, key: u8) -> char {
    if !c.is_alphanumeric() {
        return c;
    }
    let before = get_index(c);
    let after = if before >= key {
        before - key
    } else {
        before + ALPHABET_LEN - key
    };
    ALPHABET[after as usize]
}

fn get_index(c: char) -> u8 {
    if c.is_uppercase() {
        return (c as u32 - 'A' as u32) as u8 + UPPER_BASE_INDEX;
    }
    if c.is_lowercase() {
        return (c as u32 - 'a' as u32) as u8 + LOWER_BASE_INDEX;
    }
    if c.is_numeric() {
        return (c as u32 - '0' as u32) as u8 + DIGIT_BASE_INDEX;
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
        assert_eq!(
            cipher("THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG".to_string(), key),
            "qeb nrfZh Yoltk clu grjmp lsbo qeb iXwv ald"
        );
        assert_eq!(
            cipher("The password is 12345".to_string(), key),
            "q41 CxFFJBE0 5F OPQRS"
        );
        assert_eq!(
            cipher("Faça o que eu digo mas não faça o que eu faço".to_string(), key),
            "cxzx B DH1 1H 053B 9xF AxB 2xzx B DH1 1H 2xzB"
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
        assert_eq!(
            decipher("q41 CxFFJBE0 5F OPQRS".to_string(), key),
            "The password is 12345"
        );
        assert_eq!(
            decipher("cxzx B DH1 1H 053B 9xF AxB 2xzx B DH1 1H 2xzB".to_string(), key),
            "Faca o que eu digo mas nao faca o que eu faco"
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
