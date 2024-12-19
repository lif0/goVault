package unit

import "errors"

// ParseSize parses a size string (e.g., "10GB", "512MB", "1024KB", "100B") and converts it into its size in bytes.
// It accepts sizes in the following units: GB, MB, KB, and B (case-insensitive).
// The function returns an integer representing the size in bytes or an error if the input format is invalid.
// Valid input examples: "10GB", "512Mb", "1024kb", "100b".
// Invalid inputs: "abc", "10TB", "", or strings that start with non-numeric characters.
func ParseDigitalStorage(text string) (int, error) {
	if len(text) == 0 || text[0] < '0' || text[0] > '9' {
		return 0, errors.New("incorrect size")
	}

	idx := 0
	size := 0
	for idx < len(text) && text[idx] >= '0' && text[idx] <= '9' {
		number := int(text[idx] - '0')
		size = size*10 + number
		idx++
	}

	parameter := text[idx:]
	switch parameter {
	case "GB", "Gb", "gb":
		return size << 30, nil
	case "MB", "Mb", "mb":
		return size << 20, nil
	case "KB", "Kb", "kb":
		return size << 10, nil
	case "B", "b", "":
		return size, nil
	default:
		return 0, errors.New("incorrect size")
	}
}
