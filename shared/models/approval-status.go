package models

type ApprovalStatusEnum string

var (
	ApprovalStatusApproval ApprovalStatusEnum = "APPROVED"
	ApprovalStatusRejected ApprovalStatusEnum = "REJECTED"
)
