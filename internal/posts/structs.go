package posts

type Post struct {
	ID    uint64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Owner string `json:"owner" gorm:"<-:create"`
	Body  string `json:"body"`
}
