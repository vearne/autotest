package model

const (
	StateNotExecuted State = 0
	StateSuccessFul  State = 1
	StateFailed      State = 2
)

const (
	ReasonSuccess                   Reason = 0
	ReasonRequestFailed             Reason = 1
	ReasonRuleVerifyFailed          Reason = 2
	ReasonDependentItemNotCompleted Reason = 3
	ReasonTemplateRenderError       Reason = 4
	ReasonDependentItemFailed       Reason = 5
)

const (
	TypeInteger = "integer"
	TypeString  = "string"
	TypeFloat   = "float"
)
