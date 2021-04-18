package batch

// BatchItems struct defines the each batch item payload with an Id and relates to an overall BatchNo
type BatchItems struct {
	Id      int
	BatchNo int
	Item    interface{}
}
