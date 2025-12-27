package req

type IdUri struct {
	Id string `uri:"id" binding:"required,uuid"`
}

type NumberIdUri struct {
	Id int64 `uri:"id" binding:"required"`
}
