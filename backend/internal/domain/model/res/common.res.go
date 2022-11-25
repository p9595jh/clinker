package res

type SaveTxHashRes struct {
	TxHash string `json:"txHash" example:"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"`
}

type ProfuseRes[T any] struct {
	TotalCount int64 `json:"totalCount"`
	Data       []T   `json:"data"`
}
