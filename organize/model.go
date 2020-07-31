package organize

import (
	"github.com/inetspa/oneplatform-sdk-go/identity"
	uuid "github.com/satori/go.uuid"
)

type OrgClient struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"type"`
	ApiEndpoint  string `json:"api_endpoint"`
}

type OrgApiResult struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"errorMessage"`
	Code   int         `json:"code"`
}

type Department struct {
	Id           uuid.UUID          `json:"id"`
	Name         string             `json:"dept_name"`
	ParentDeptId *uuid.UUID         `json:"parent_dept_id"`
	Accounts     *identity.Employee `json:"has_account"`
}

type TeamMember struct {
	Id           uuid.UUID            `json:"id"`
	Name         string               `json:"dept_name"`
	ParentDeptId *uuid.UUID           `json:"parent_dept_id"`
	Accounts     *[]identity.Employee `json:"has_account"`
}

type HeadDepartment struct {
	Id           uuid.UUID            `json:"id"`
	Name         string               `json:"dept_name"`
	ParentDeptId *uuid.UUID           `json:"parent_dept_id"`
	Accounts     *[]identity.Employee `json:"has_account"`
}
