package models

type ReportStatusEnum string

var (
	ReportStatusNew        ReportStatusEnum = "NEW"
	ReportStatusEvaluating ReportStatusEnum = "EVALUATING"
	ReportStatusPending    ReportStatusEnum = "PENDING"
	ReportStatusPublished  ReportStatusEnum = "PUBLISHED"
	ReportStatusError      ReportStatusEnum = "ERROR"
)
