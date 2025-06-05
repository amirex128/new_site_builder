package domain

type Configuration struct {
	ID    int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Key   string `json:"key" gorm:"column:key;type:varchar(255);unique;not null"`
	Value string `json:"value" gorm:"column:value;type:text;not null"`
}

func (Configuration) TableName() string {
	return "configurations"
}
