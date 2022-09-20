package dto

// Example
type Tenant struct {
}

type TenantRequestV1 struct {
	Name       string `form:"name" json:"name" binding:"required"`
	SeparateDb bool   `form:"separateDb" json:"separateDb"`
}
