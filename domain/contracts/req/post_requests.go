package req

type AddPost struct {
	Url string `json:"url" binding:"required"`
}
