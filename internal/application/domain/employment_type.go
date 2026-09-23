package domain

type EmploymentType string

const (
	StatusEmployed     EmploymentType = "EMPLOYED"
	StatusSelfEmployed EmploymentType = "SELF_EMPLOYED"
	StatusUnemployed   EmploymentType = "UNEMPLOYED"
	StatusStudent      EmploymentType = "STUDENT"
	StatusRetired      EmploymentType = "RETIRED"
)
