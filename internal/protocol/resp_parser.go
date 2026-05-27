package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if len(line) == 0 || line[0] != '*' {
		return nil, errors.New("invalid RESP array")
	}

	count, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}

	parts := make([]string, 0, count)

	for i := 0; i < count; i++ {
		bulkHeader, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if len(bulkHeader) == 0 || bulkHeader[0] != '$' {
			return nil, fmt.Errorf("invalid RESP bulk string")
		}

		bulkLen, err := strconv.Atoi(strings.TrimSpace(bulkHeader[1:]))
		if err != nil {
			return nil, err
		}

		if bulkLen == -1 {
			parts = append(parts, "")
			continue
		}

		buf := make([]byte, bulkLen+2)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, err
		}

		parts = append(parts, string(buf[:bulkLen]))
	}

	return parts, nil
}
