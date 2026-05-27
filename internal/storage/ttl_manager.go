package storage

import "time"

func StartTTLWorker() {

	go func() {

		for {

			time.Sleep(time.Second)

			Manager.ForEachDB(func(db *DB) {
				db.Store.RemoveExpiredKeys()
			})
		}
	}()
}
