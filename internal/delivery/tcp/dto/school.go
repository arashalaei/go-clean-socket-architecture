package dto

type School struct {
	Id      uint    `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Classes []Class `json:"classes,omitempty"`
}
