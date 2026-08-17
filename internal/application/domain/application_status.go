package domain

type ApplicationStatus string

const (
	StatusNew        ApplicationStatus = "NEW"
	StatusValidating ApplicationStatus = "VALIDATING"
	StatusScoring    ApplicationStatus = "SCORING"
	StatusApproved   ApplicationStatus = "APPROVED"
	StatusRejected   ApplicationStatus = "REJECTED"
)
