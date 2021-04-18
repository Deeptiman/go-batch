package batch

// BatchItems struct defines the each batch item payload with an Id and relates to an overall BatchNo
type BatchItems struct {
	Id      int         `json:"id"`
	BatchNo int         `json:"batchNo"`
	Item    interface{} `json:"item"`
}
