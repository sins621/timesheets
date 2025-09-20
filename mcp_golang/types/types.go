package types

import "time"

type Project struct {
	TaskID    int       `json:"TaskId"`
	Name      string    `json:"Name"`
	IsActive  bool      `json:"IsActive"`
	CreatedOn time.Time `json:"Created_On"`
	UpdatedOn time.Time `json:"Updated_On"`
	Client    struct {
		GroupID  int    `json:"GroupId"`
		Name     string `json:"Name"`
		Currency string `json:"Currency"`
	} `json:"Client"`
}

type WorkEntry struct {
	Description string
	Date        time.Time
	Hours       int
	TaskID      int
	CostCodeID  int
}
