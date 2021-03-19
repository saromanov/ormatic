package models

// Create defines model for create table
type Create struct {
	TableName   string
	TableFields []TableField
}

// TableField defines field for the table
type TableField struct {
	Name string
	Type string
	Tags Tags
}
