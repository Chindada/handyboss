package models

// TableName TableName
func (c *Wise) TableName() string {
	return "fake_data_wise"
}

// TableName TableName
func (c *Di) TableName() string {
	return "fake_data_di"
}

// TableIndex TableIndex
func (c *Di) TableIndex() [][]string {
	return [][]string{
		{
			"MacAddress",
			"Timestamp",
		},
	}
}

// TableName TableName
func (c *DcStatus) TableName() string {
	return "fake_data_status"
}

// TableIndex TableIndex
func (c *DcStatus) TableIndex() [][]string {
	return [][]string{
		{
			"MacAddress",
			"Timestamp",
		},
	}
}
