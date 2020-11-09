package models

// Resource RealStatus
type Resource struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Icon      string  `json:"icon"`
	Seq       int64   `json:"seq"`
	Operation []int64 `json:"operation"`
}

// LoginReturn StatusReturn
type LoginReturn struct {
	Response string     `json:"response"`
	Resource []Resource `json:"resource"`
	User     User       `json:"user"`
	Edition  string     `json:"edition"`
}

// User User
type User struct {
	Data        []RealStatus `json:"data"`
	Response    string       `json:"response"`
	ID          int64        `json:"id"`
	Account     string       `json:"account"`
	Password    string       `json:"password"`
	Name        string       `json:"name"`
	Status      bool         `json:"status"`
	PhoneNumber string       `json:"phoneNumber"`
	Email       string       `json:"email"`
	Avatar      string       `json:"avatar"`
	RoleID      int64        `json:"roleId"`
}
