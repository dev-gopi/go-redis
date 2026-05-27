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
