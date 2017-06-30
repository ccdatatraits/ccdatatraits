package main

import (
	"fmt"
)

const (
	//ENDPOINT = ""
)

func main() {
	// proxy?

	sample := `{
        "key": 22797,
        "offset": 16018,
        "partition": 2,
        "value": {
            "changes": {
                "history": {
                    "com.mediamath.changes.change": {
                        "api_key": null,
                        "entity_id": 22797,
                        "entity_type_id": 32,
                        "is_delete": null,
                        "snapshot": {
                            "com.mediamath.changes.snapshot": {
                                "last_snapshotted": 398489318,
                                "last_topic": "ewr.changes.history.t1db.1",
                                "snapshot_time": 1490302362,
                                "snapshot_type": "BOOTSTRAP",
                                "snapshotter_api_key": "qn269kmekkhg8xh4f56uk5aj",
                                "snapshotter_sha": "1f93307"
                            }
                        },
                        "source_time": null,
                        "user_id": null
                    }
                },
                "publish_time": 1490302356,
                "publisher_api_key": "qn269kmekkhg8xh4f56uk5aj"
            },
            "fields": {
                "created_on": "2015-02-23 00:43:03.891647+00",
                "currency_code": "PLN",
                "date": "2015-02-23",
                "id": "22797",
                "rate": "3.6650426557",
                "updated_on": "2015-02-23 00:43:03.891647+00",
                "version": "0"
            },
            "lsn": null,
            "publisher_sha": "1f93307",
            "wal_buffer_offset": "\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000",
            "wal_name": "bootstrapper",
            "xid": null
        }
    }
	`
	fmt.Println(sample)
}
