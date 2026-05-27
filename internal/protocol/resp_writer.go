package protocol

import "fmt"

func SimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

func BulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func NullBulkString() string {
	return "$-1\r\n"
}

func Error(s string) string {
	return fmt.Sprintf("-%s\r\n", s)
}

func Integer(i int) string {
	return fmt.Sprintf(":%d\r\n", i)
}

func Array(values []string) string {

	resp := "*" + fmt.Sprintf("%d", len(values)) + "\r\n"

	for _, v := range values {
		resp += BulkString(v)
	}

	return resp
}

func ArrayWithNulls(values []string, present []bool) string {

	resp := "*" + fmt.Sprintf("%d", len(values)) + "\r\n"

	for i, v := range values {
		if i < len(present) && !present[i] {
			resp += NullBulkString()
			continue
		}

		resp += BulkString(v)
	}

	return resp
}
