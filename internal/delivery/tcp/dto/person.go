package dto

type Person struct {
	Id      uint   `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Classes []uint `json:"calasses,omitempty"`
}
