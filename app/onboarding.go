package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"onboarding-process/models"
	"time"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/workflow"
)

var (
	onboardingStateStoreName    = "statestore"
	onboardingWorkflowComponent = "dapr"
	onboardingWorkflowName      = "ProcessWorkflow"
	onboardingDefaultName       = "TestOrg01"
)

/*
## Demo App Process (Onboard a new organization, create and add users to it, and add users to existing groups)
- Initialise Dapr Client.
- Create a new Organization.
- Create 2 x new users.
- Add users to existing groups.
*/
func OnboardNewOrganization() {
	fmt.Println("*** Welcome to the Dapr Workflow console demo!")
	fmt.Printf("*** This app (%v - %v), will place a request that will start the onboarding workflow process.\n", onboardingWorkflowComponent, onboardingWorkflowName)

	fmt.Println("==========Begin the initialisation process:==========")

	fmt.Println("*** Initialising the Dapr Worker Client.")
	w, err := workflow.NewWorker()
	if err != nil {
		log.Fatalf("failed to start worker: %v", err)
	}
	fmt.Println("*** Completed initialising the Dapr Worker Client.")

	fmt.Println("*** Using the worker to register workflow activities (RegisterWorkflow adds a workflow function to the registry).")
	if err := w.RegisterWorkflow(ProcessWorkflow); err != nil {
		log.Fatal(err)
	}
	if err := w.RegisterActivity(NotifyActivity); err != nil {
		log.Fatal(err)
	}
	if err := w.RegisterActivity(RequestApprovalActivity); err != nil {
		log.Fatal(err)
	}
	if err := w.RegisterActivity(VerifyOnboardingActivity); err != nil {
		log.Fatal(err)
	}
	if err := w.RegisterActivity(ProcessAddUserActivity); err != nil {
		log.Fatal(err)
	}
	if err := w.RegisterActivity(UpdateOnboardingActivity); err != nil {
		log.Fatal(err)
	}
	fmt.Println("*** Completed registering workflow activities.")

	fmt.Println("*** Initialise a non-blocking worker to handle workflows and activities registered prior to this being called.")
	if err := w.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("*** Completed initialisation.")

	fmt.Println("*** Initialise Dapr Client and wfClient.")
	daprClient, err := client.NewClient()
	if err != nil {
		log.Fatalf("failed to initialise dapr client: %v", err)
	}
	wfClient, err := workflow.NewClient(workflow.WithDaprClient(daprClient))
	if err != nil {
		log.Fatalf("failed to initialise workflow client: %v", err)
	}
	fmt.Println("*** Completed initialising Dapr Client and wfClient.")

	onboardingRequest := models.WorkflowItem{ItemName: "TestOrg01", NumOfUsers: 3}
	if err := createNewOrganizationRequest(daprClient, onboardingRequest); err != nil {
		log.Fatalf("failed to log new organization request: %v", err)
	}
	fmt.Println("==========End of the initialisation process:==========")

	fmt.Println("==========Begin the onboarding process:==========")

	itemName := onboardingDefaultName
	numberOfUsers := onboardingRequest.NumOfUsers

	onboardingPayload := models.OnboardingPayload{
		ItemName:   itemName,
		NumOfUsers: numberOfUsers,
	}

	id, err := wfClient.ScheduleNewWorkflow(context.Background(), onboardingWorkflowName, workflow.WithInput(onboardingPayload))
	if err != nil {
		log.Fatalf("failed to start workflow: %v", err)
	}

	approvalSought := false

	startTime := time.Now()

	for {
		timeDelta := time.Since(startTime)
		metadata, err := wfClient.FetchWorkflowMetadata(context.Background(), id)
		if err != nil {
			log.Fatalf("failed to fetch workflow: %v", err)
		}
		if (metadata.RuntimeStatus == workflow.StatusCompleted) || (metadata.RuntimeStatus == workflow.StatusFailed) || (metadata.RuntimeStatus == workflow.StatusTerminated) {
			fmt.Printf("Workflow completed - result: %v\n", metadata.RuntimeStatus.String())
			break
		}
		if timeDelta.Seconds() >= 10 {
			metadata, err := wfClient.FetchWorkflowMetadata(context.Background(), id)
			if err != nil {
				log.Fatalf("failed to fetch workflow: %v", err)
			}
			runtimeStatusHasNotCompleted := metadata.RuntimeStatus != workflow.StatusCompleted
			runtimeStatusHasNotFailed := metadata.RuntimeStatus != workflow.StatusFailed
			if numberOfUsers > 50 && !approvalSought && (runtimeStatusHasNotCompleted || runtimeStatusHasNotFailed || (metadata.RuntimeStatus != workflow.StatusTerminated)) {
				approvalSought = true
				promptForApprovalRequest(id)
			}
		}
		// Sleep before the next iteration
		time.Sleep(time.Second)
	}

	fmt.Println("Onboarding of organization is complete")
}

// promptForApproval is an example case. There is no user input required here due to this being for testing purposes only.
// It would be perfectly valid to add a wait here or display a prompt to continue the process.
func promptForApprovalRequest(id string) {
	wfClient, err := workflow.NewClient()
	if err != nil {
		log.Fatalf("failed to initialise wfClient: %v", err)
	}
	if err := wfClient.RaiseEvent(context.Background(), id, "manager_approval"); err != nil {
		log.Fatal(err)
	}
}

func createNewOrganizationRequest(daprClient client.Client, org models.WorkflowItem) error {
	item := org
	itemSerialized, err := json.Marshal(item)
	if err != nil {
		return err
	}
	fmt.Printf("creating request for new organization entity: %s\n", item.ItemName)
	if err := daprClient.SaveState(context.Background(), onboardingStateStoreName, item.ItemName, itemSerialized, nil); err != nil {
		return err
	}
	return nil
}
