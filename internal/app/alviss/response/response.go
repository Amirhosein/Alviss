package response

type Message struct {
	Message string `json:"message" binding:"required"`
}

type SuccessfullyCreated struct {
	Message  string `json:"message" binding:"required"`
	ShortURL string `json:"ShortURL" binding:"required"`
}

type Detail struct {
	OriginalURL string `json:"OriginalURL" binding:"required"`
	ShortURL    string `json:"ShortURL" binding:"required"`
	UsedCount   int    `json:"UsedCount" binding:"required"`
	ExpDate     string `json:"ExpDate" binding:"required"`
}
