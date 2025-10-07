package idgen

// alphabet is a scrambled base62 alphabet for unpredictable encoding.
var alphabet = []rune{
	'z', 'K', '7', 'p', 'B', 'x', 'N', '4', 'i', 'S',
	'9', 'u', 'E', 'R', 'c', '0', 'm', 'W', 'I', '2',
	'y', 'L', '6', 'q', 'A', 'w', 'M', '5', 'j', 'T',
	'8', 'v', 'F', 'Q', 'd', '1', 'n', 'X', 'J', '3',
	'Z', 'O', 'e', 'r', 'C', 'Y', 'k', 'U', 'G', 't',
	'f', 'V', 'P', 'D', 's', 'H', 'a', 'b', 'o', 'g',
	'h', 'l',
}

// Encode converts an integer ID to a base62 string using the scrambled alphabet.
func Encode(id int64) string {
	if id == 0 {
		return string(alphabet[0])
	}

	var digits []rune
	for id > 0 {
		digits = append([]rune{alphabet[id%62]}, digits...)
		id /= 62
	}

	return string(digits)
}
