package domain

// Permission represents User.Permissions table
type Permission struct {
	ID   int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement;type:bigint"`
	Name string `json:"name" gorm:"column:name;type:longtext;not null"`

	// Relations
	Roles []Role `json:"roles" gorm:"many2many:permission_roles;"`
}

// TableName specifies the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}
