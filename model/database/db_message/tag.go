package db_message

type Tag struct {
	Id  uint   `json:"id" gorm:"primary_key;column:id"`
	Tag string `json:"tag"`
}

func (Tag) TableName() string {
	return "tag_tab"
}
