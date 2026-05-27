package commands

import (
	"sort"
	"strconv"
	"strings"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

type zRangeItem struct {
	member string
	score  float64
}

func HandleZRange(cl *client.Client, cmd []string) string {
	if len(cmd) < 4 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	start, err := strconv.Atoi(cmd[2])
	if err != nil {
		return protocol.Error("invalid start index")
	}

	stop, err := strconv.Atoi(cmd[3])
	if err != nil {
		return protocol.Error("invalid stop index")
	}

	withScores := false
	for i := 4; i < len(cmd); i++ {
		if strings.EqualFold(cmd[i], "WITHSCORES") {
			withScores = true
		}
	}

	db := storage.GetClientDB(cl)
	value, exists := db.Store.GetValue(key)
	if !exists {
		return protocol.Array([]string{})
	}
	if value.Type != storage.ZSetType {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	zset, ok := db.Store.GetZSet(key)
	if !ok || len(zset) == 0 {
		return protocol.Array([]string{})
	}

	items := make([]zRangeItem, 0, len(zset))
	for member, score := range zset {
		items = append(items, zRangeItem{member: member, score: score})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].score == items[j].score {
			return items[i].member < items[j].member
		}
		return items[i].score < items[j].score
	})

	length := len(items)
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start >= length || start > stop {
		return protocol.Array([]string{})
	}

	result := make([]string, 0, (stop-start+1)*(1+boolToInt(withScores)))
	for i := start; i <= stop; i++ {
		result = append(result, items[i].member)
		if withScores {
			result = append(result, strconv.FormatFloat(items[i].score, 'f', -1, 64))
		}
	}

	return protocol.Array(result)
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
