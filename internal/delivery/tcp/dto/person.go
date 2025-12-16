package dto

type Person struct {
	Id      uint   `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Role    string `json:"role,omitempty"`
	SchoolId uint  `json:"school_id,omitempty"`
	Classes []uint `json:"classes,omitempty"`
}

type CreatePersonReq struct {
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
	SchoolId uint   `json:"school_id,omitempty"`
}

type WhoAmIReq struct {
	PersonId uint `json:"person_id,omitempty"`
}
