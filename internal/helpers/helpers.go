package helpers

import (
	"encoding/base64"
	"fmt"
)

// encodeCursor converts an offset into a base64 cursor
func EncodeCursor(offset int32) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("offset:%d", offset)))
}

// decodeCursor converts a base64 cursor back into an offset
func DecodeCursor(cursor string) (int32, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}

	var offset int32
	_, err = fmt.Sscanf(string(decoded), "offset:%d", &offset)
	if err != nil {
		return 0, err
	}

	return offset, nil
}
