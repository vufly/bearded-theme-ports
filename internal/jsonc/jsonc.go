package jsonc

import (
	"encoding/json"
	"os"
)

func UnmarshalFile(path string, value any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return Unmarshal(content, value)
}

func Unmarshal(content []byte, value any) error {
	return json.Unmarshal(stripLineComments(content), value)
}

func stripLineComments(content []byte) []byte {
	clean := make([]byte, 0, len(content))
	inString := false
	escaped := false
	inComment := false

	for index := 0; index < len(content); index++ {
		char := content[index]

		if inComment {
			if char == '\n' {
				inComment = false
				clean = append(clean, char)
			}
			continue
		}

		if inString {
			clean = append(clean, char)
			if escaped {
				escaped = false
				continue
			}
			if char == '\\' {
				escaped = true
				continue
			}
			if char == '"' {
				inString = false
			}
			continue
		}

		if char == '"' {
			inString = true
			clean = append(clean, char)
			continue
		}

		if char == '/' && index+1 < len(content) && content[index+1] == '/' {
			inComment = true
			index++
			continue
		}

		clean = append(clean, char)
	}

	return clean
}
