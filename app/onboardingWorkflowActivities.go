package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"onboarding-process/models"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/workflow"
)

// NotifyActivity outputs a notification message
func NotifyActivity(ctx workflow.ActivityContext) (any, error) {
	var input models.Notification
	if err := ctx.GetInput(&input); err != nil {
		return "", err
	}
	fmt.Printf("NotifyActivity: %s\n", input.Message)
	return nil, nil
}

// ProcessAddUserActivity is used to process adding users to the organization
func ProcessAddUserActivity(ctx workflow.ActivityContext) (any, error) {
	var input models.OnboardingRequest

	if err := ctx.GetInput(&input); err != nil {
		return "", err
	}

	fmt.Printf("ProcessAddUserActivity: Processing request %s adding %d users to organization %s\n", input.RequestID, input.NumOfUsers, input.ItemBeingProcessed)
	return nil, nil
}

// VerifyOnboardingActivity is used to verify if an onboarding request has been scheduled
func VerifyOnboardingActivity(ctx workflow.ActivityContext) (any, error) {
	var input models.WorkflowRequest
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}

	fmt.Printf("VerifyOnboardingActivity: Verifying request for %s onboarding %s with %d users\n", input.RequestID, input.RequestName, input.NumOfUsers)
	dClient, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	item, err := dClient.GetState(context.Background(), onboardingStateStoreName, input.RequestName, nil)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return models.WorkflowResult{
			Success:      false,
			WorkflowItem: models.WorkflowItem{},
		}, nil
	}

	var result models.WorkflowItem
	if err := json.Unmarshal(item.Value, &result); err != nil {
		log.Fatalf("failed to parse workflow result %v", err)
	}
	fmt.Printf("VerifyOnboardingActivity: The onboarding request for %s with %d user/s is ready for processing\n", result.ItemName, result.NumOfUsers)
	if result.NumOfUsers >= input.NumOfUsers {
		return models.WorkflowResult{Success: true, WorkflowItem: result}, nil
	}
	return models.WorkflowResult{Success: false, WorkflowItem: models.WorkflowItem{}}, nil
}

// UpdateOnboardingActivity modifies the workflow logs.
func UpdateOnboardingActivity(ctx workflow.ActivityContext) (any, error) {
	var input models.OnboardingRequest
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}
	fmt.Printf("UpdateOnboardingActivity: Checking Request %s for onboarding %s with %d user/s\n", input.RequestID, input.ItemBeingProcessed, input.NumOfUsers)
	dClient, err := client.NewClient()
	if err != nil {
		return nil, err
	}
	item, err := dClient.GetState(context.Background(), onboardingStateStoreName, url.QueryEscape(input.ItemBeingProcessed), nil)
	if err != nil {
		return nil, err
	}

	var result models.WorkflowItem
	err = json.Unmarshal(item.Value, &result)
	if err != nil {
		return nil, err
	}
	newQuantity := result.NumOfUsers - input.NumOfUsers
	if newQuantity < 0 {
		return nil, fmt.Errorf("incorrect number of organisations being onboarding for: %s", input.ItemBeingProcessed)
	}

	result.NumOfUsers = input.NumOfUsers
	newState, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("failed to marshal new state: %v", err)
	}
	dClient.SaveState(context.Background(), onboardingStateStoreName, input.ItemBeingProcessed, newState, nil)
	fmt.Printf("UpdateOnboardingActivity: The request for new organization %s with %d user/s have been successfully validated\n", result.ItemName, result.NumOfUsers)
	return models.WorkflowResult{Success: true, WorkflowItem: result}, nil
}

// RequestApprovalActivity requests approval for the order
func RequestApprovalActivity(ctx workflow.ActivityContext) (any, error) {
	var input models.OnboardingPayload
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}
	fmt.Printf("RequestApprovalActivity: Requesting approval for onboarding new organisation %s with %d user/s\n", input.ItemName, input.NumOfUsers)
	return models.ApprovalRequired{Approval: true}, nil
}
