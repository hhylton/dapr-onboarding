package models

type OnboardingPayload struct {
	ItemName   string `json:"item_name"`
	NumOfUsers int    `json:"numofusers"`
}

type OnboardingResult struct {
	Processed bool `json:"processed"`
}

type WorkflowItem struct {
	ItemName   string `json:"workflow_name"`
	NumOfUsers int    `json:"numofusers"`
}

type WorkflowRequest struct {
	RequestID   string `json:"request_id"`
	RequestName string `json:"request_name"`
	NumOfUsers  int    `json:"numofusers"`
}

type WorkflowResult struct {
	Success      bool         `json:"success"`
	WorkflowItem WorkflowItem `json:"workflow_item"`
}

type OnboardingRequest struct {
	RequestID          string `json:"request_id"`
	ItemBeingProcessed string `json:"item_being_processed"`
	NumOfUsers         int    `json:"numofusers"`
}

type ApprovalRequired struct {
	Approval bool `json:"approval"`
}

type Notification struct {
	Message string `json:"message"`
}
