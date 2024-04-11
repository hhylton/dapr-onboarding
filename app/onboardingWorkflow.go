package app

import (
	"fmt"
	"log"
	"onboarding-process/models"
	"time"

	"github.com/dapr/go-sdk/workflow"
)

// ProcessWorkflow is the main workflow for orchestrating activities in the onboarding process.
func ProcessWorkflow(ctx *workflow.WorkflowContext) (any, error) {
	onboardingID := ctx.InstanceID()
	var onboardingPayload models.OnboardingPayload
	if err := ctx.GetInput(&onboardingPayload); err != nil {
		return nil, err
	}

	err := ctx.CallActivity(NotifyActivity, workflow.ActivityInput(models.Notification{
		Message: fmt.Sprintf("Received onboarding request %s for new organization %s with %d users", onboardingID, onboardingPayload.ItemName, onboardingPayload.NumOfUsers),
	})).Await(nil)
	if err != nil {
		return models.OnboardingResult{Processed: false}, err
	}

	var verifyWorkflowResult models.WorkflowResult
	if err := ctx.CallActivity(VerifyOnboardingActivity, workflow.ActivityInput(models.WorkflowRequest{
		RequestID:   onboardingID,
		RequestName: onboardingPayload.ItemName,
		NumOfUsers:  onboardingPayload.NumOfUsers,
	})).Await(&verifyWorkflowResult); err != nil {
		return models.OnboardingResult{Processed: false}, err
	}

	if !verifyWorkflowResult.Success {
		notification := models.Notification{Message: fmt.Sprintf(" for %s", onboardingPayload.ItemName)}
		err := ctx.CallActivity(NotifyActivity, workflow.ActivityInput(notification)).Await(nil)
		return models.OnboardingResult{Processed: false}, err
	}

	if onboardingPayload.NumOfUsers > 50 {
		var approvalRequired models.ApprovalRequired
		if err := ctx.CallActivity(RequestApprovalActivity, workflow.ActivityInput(onboardingPayload)).Await(&approvalRequired); err != nil {
			return models.OnboardingResult{Processed: false}, err
		}
		if err := ctx.WaitForExternalEvent("manager_approval", time.Second*200).Await(nil); err != nil {
			return models.OnboardingResult{Processed: false}, err
		}
		// TODO: Confirm timeout flow - this will be in the form of an error.
		if approvalRequired.Approval {
			if err := ctx.CallActivity(NotifyActivity, workflow.ActivityInput(models.Notification{Message: fmt.Sprintf("Onboarding for request %s has been approved!", onboardingID)})).Await(nil); err != nil {
				log.Printf("failed to notify of a successful onboarding: %v\n", err)
			}
		} else {
			if err := ctx.CallActivity(NotifyActivity, workflow.ActivityInput(models.Notification{Message: fmt.Sprintf("Onboarding for request %s has been rejected!", onboardingID)})).Await(nil); err != nil {
				log.Printf("failed to notify of an unsuccessful onboarding :%v\n", err)
			}
			return models.OnboardingResult{Processed: false}, err
		}
	}
	err = ctx.CallActivity(ProcessAddUserActivity, workflow.ActivityInput(models.OnboardingRequest{
		RequestID:          onboardingID,
		ItemBeingProcessed: onboardingPayload.ItemName,
		NumOfUsers:         onboardingPayload.NumOfUsers,
	})).Await(nil)
	if err != nil {
		if err := ctx.CallActivity(NotifyActivity, workflow.ActivityInput(models.Notification{Message: fmt.Sprintf("Onboarding %s failed!", onboardingID)})).Await(nil); err != nil {
			log.Printf("failed to notify of a failed onboarding: %v", err)
		}
		return models.OnboardingResult{Processed: false}, err
	}

	err = ctx.CallActivity(UpdateOnboardingActivity, workflow.ActivityInput(models.OnboardingRequest{
		RequestID:          onboardingID,
		ItemBeingProcessed: onboardingPayload.ItemName,
		NumOfUsers:         onboardingPayload.NumOfUsers,
	})).Await(nil)
	if err != nil {
		if err := ctx.CallActivity(NotifyActivity, workflow.ActivityInput(models.Notification{Message: fmt.Sprintf("Onboarding %s failed!", onboardingID)})).Await(nil); err != nil {
			log.Printf("failed to notify of a failed onboarding: %v", err)
		}
		return models.OnboardingResult{Processed: false}, err
	}

	if err := ctx.CallActivity(NotifyActivity, workflow.ActivityInput(models.Notification{Message: fmt.Sprintf("The onboarding process %s for %s has successfully completed!", onboardingID, onboardingPayload.ItemName)})).Await(nil); err != nil {
		log.Printf("failed to notify of a successful onboarding: %v", err)
	}
	return models.OnboardingResult{Processed: true}, err
}
