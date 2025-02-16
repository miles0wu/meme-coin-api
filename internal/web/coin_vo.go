package web

type CoinVo struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updated"`
	PopularityScore uint32 `json:"popularityScore"`
}

type CreateCoinReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCoinReq struct {
	Description string `json:"description"`
}
