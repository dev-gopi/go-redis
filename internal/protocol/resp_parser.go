package protocol

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

func ParseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, errors.New("invalid RESP array")
	}

	count, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}

	parts := make([]string, 0, count)

	for i := 0; i < count; i++ {
		_, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		data, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		parts = append(parts, strings.TrimSpace(data))
	}

	return parts, nil
}
