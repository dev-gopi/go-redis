package aof

import (
	"os"
	"time"

	"github.com/dev-gopi/go-redis/internal/logger"
)

func StartAutoRotate(
	path string,
	maxSize int64,
) {

	go func() {

		for {

			time.Sleep(time.Minute)

			info, err := os.Stat(path)
			if err != nil {
				continue
			}

			if info.Size() < maxSize {
				continue
			}

			backup := path + "." +
				time.Now().
					Format("20060102150405")

			err = os.Rename(path, backup)
			if err != nil {
				continue
			}

			logger.InfoLogger.Println(
				"AOF rotated:",
				backup,
			)

			Init(path)
		}
	}()
}
