package controllers

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CreateAllPostsInput struct {
	Posts []CreatePostInput `json:"posts"`
}

type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
