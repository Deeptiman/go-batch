package batch

import "encoding/json"

func getJsonString(item interface{}) string {
	batchItem, _ := json.Marshal(item)
	return string(batchItem)
}
