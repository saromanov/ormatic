package models

// Create defines model for create table
type Create struct {
	TableName     string
	TableFields   []TableField
	Relationships []Relationship
}

// TableField defines field for the table
type TableField struct {
	Name string
	Type string
	Tags Tags
}

// Relationship is struct for defining relationships at betweeb tables
type Relationship struct {
	Parent    string
	TableName string
	Name      string
	Column    string
}
