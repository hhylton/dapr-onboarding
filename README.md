# Dapr Onboarding Workflow

## Introduction to Dapr

Dapr is a portable, event-driven runtime that makes it easy for any developer to build resilient, stateless, and stateful applications that run on the cloud and edge and embraces the diversity of languages and developer frameworks.

With the current trend for developing microservice application architectures, which are inherently distributed. Dapr codifies the best practices for building microservice applications into open, independent APIs called building blocks. Dapr’s building blocks:

- Enable you to build portable applications using the language and framework of your choice.
- Are completely independent
- Have no limit to how many you use in your application

Dapr is platform agnostic, meaning you can run your applications:

- Locally
- On any Kubernetes cluster
- On virtual or physical machines
- In other hosting environments that Dapr integrates with.
- This enables you to build microservice applications that can run on the cloud and edge.

### Building blocks <span style="font-size: 0.5em">[link](https://docs.dapr.io/developing-applications/building-blocks/)</span>

Dapr capabilities that solve common development challenges for distributed applications:

- [Service invocation](https://docs.dapr.io/developing-applications/building-blocks/service-invocation/): Perform direct, secure, service-to-service method calls
- [State management](https://docs.dapr.io/developing-applications/building-blocks/state-management/): Create long running stateful services
- [Publish & subscribe messaging](https://docs.dapr.io/developing-applications/building-blocks/pubsub/): Secure, scalable messaging between services
- [Bindings](https://docs.dapr.io/developing-applications/building-blocks/bindings/): Interface with or be triggered from external systems
- [Actors](https://docs.dapr.io/developing-applications/building-blocks/actors/): Encapsulate code and data in reusable actor objects as a common microservices design pattern
- [Secrets management](https://docs.dapr.io/developing-applications/building-blocks/secrets/): Securely access secrets from your application
- [Configuration](https://docs.dapr.io/developing-applications/building-blocks/configuration/): Manage and be notified of application configuration changes
- [Distributed lock](https://docs.dapr.io/developing-applications/building-blocks/distributed-lock/): Distributed locks provide mutually exclusive access to shared resources from an application.
- [Workflow](https://docs.dapr.io/developing-applications/building-blocks/workflow/): Orchestrate logic across various microservices
- [Cryptography](https://docs.dapr.io/developing-applications/building-blocks/cryptography/): Perform cryptographic operations without exposing keys to your application

## The Demo App

The demo is a simple console application to demonstrate Dapr's workflow programming model and the workflow authoring client. The console app 
starts and manages an organization onboarding process workflow.

This application includes one project:

- Go app `dapr-onboarding` 

The application contains 1 workflow (ProcessWorkflow) which simulates the process of onboarding a new organization, and 5 unique activities within the workflow. These 5 activities are as follows:

- NotifyActivity: This activity utilizes a logger to print out messages throughout the workflow. These messages notify the user when the onboarding request has been scheduled, also when an error occurs, and more.
- AddUserActivity: This activity is used to process creating users and adding them to the organization. It also assigns users to existing groups.
- VerifyOnboardingActivity: This activity is used to verify if an onboarding request has been correctly scheduled. It does this by checking the state store to ensure that the request has been correctly logged.
- UpdateOnboardingActivity: This activity updates the state store to flag the onbording request as complete.
- RequestApprovalActivity: This activity requests approval from a manager, if the organisation has more than 50 users.

### Run the onboarding process workflow

1. Open a new terminal window and navigate to root of the project, the  `dapr-onboarding` directory.
2. Run the console app with Dapr:

  ```sh
  dapr run -f .
  ```

3. Expected output:

  ```sh
ℹ️  Validating config and starting app "onboarding-process"
ℹ️  Started Dapr with app id "onboarding-process". HTTP Port: 43011. gRPC Port: 37599
ℹ️  Writing log files to directory : /home/hhyl/proj/github/workflow-orchestration/dapr/dapr-onboarding/cmd/.dapr/logs
== APP - onboarding-process == *** Welcome to the Dapr Workflow console demo!
== APP - onboarding-process == *** This app (dapr - ProcessWorkflow), will place a request that will start the onboarding workflow process.
== APP - onboarding-process == ==========Begin the initialisation process:==========
== APP - onboarding-process == *** Initialising the Dapr Worker Client.
== APP - onboarding-process == dapr client initializing for: 127.0.0.1:37599
== APP - onboarding-process == *** Completed initialising the Dapr Worker Client.
== APP - onboarding-process == *** Using the worker to register workflow activities (RegisterWorkflow adds a workflow function to the registry).
== APP - onboarding-process == *** Completed registering workflow activities.
== APP - onboarding-process == *** Initialise a non-blocking worker to handle workflows and activities registered prior to this being called.
== APP - onboarding-process == *** Completed initialisation.
== APP - onboarding-process == 2024/04/11 09:18:20 work item listener started
== APP - onboarding-process == *** Initialise Dapr Client and wfClient.
== APP - onboarding-process == INFO: 2024/04/11 09:18:20 starting background processor
== APP - onboarding-process == *** Completed initialising Dapr Client and wfClient.
== APP - onboarding-process == creating request for new organization entity: TestOrg01
== APP - onboarding-process == ==========End of the initialisation process:==========
== APP - onboarding-process == ==========Begin the onboarding process:==========
== APP - onboarding-process == NotifyActivity: Received onboarding request ceb60a94-e598-4b2f-af7f-8ac7e5c683ec for new organization TestOrg01 with 3 users
== APP - onboarding-process == VerifyOnboardingActivity: Verifying request for ceb60a94-e598-4b2f-af7f-8ac7e5c683ec onboarding TestOrg01 with 3 users
== APP - onboarding-process == VerifyOnboardingActivity: The onboarding request for TestOrg01 with 3 user/s is ready for processing
== APP - onboarding-process == ProcessAddUserActivity: Processing request ceb60a94-e598-4b2f-af7f-8ac7e5c683ec adding 3 users to organization TestOrg01
== APP - onboarding-process == UpdateOnboardingActivity: Checking Request ceb60a94-e598-4b2f-af7f-8ac7e5c683ec for onboarding TestOrg01 with 3 user/s
== APP - onboarding-process == UpdateOnboardingActivity: The request for new organization TestOrg01 with 3 user/s have been successfully validated
== APP - onboarding-process == NotifyActivity: The onboarding process ceb60a94-e598-4b2f-af7f-8ac7e5c683ec for TestOrg01 has successfully completed!
== APP - onboarding-process == Workflow completed - result: COMPLETED
== APP - onboarding-process == Onboarding of organization is complete
Exited App successfully
  ```

4. Stop Dapr workflow with CTRL-C or:

```sh
dapr stop -f .
```

### View workflow output with Zipkin

For a more detailed view of the workflow activities (duration, progress etc.), try using Zipkin.

1. View Traces in Zipkin UI - In your browser go to http://localhost:9411 to view the workflow trace spans in the Zipkin web UI. The onboarding-process workflow should be viewable with the following output in the Zipkin web UI. Note: the [openzipkin/zipkin](https://hub.docker.com/r/openzipkin/zipkin/) docker container is 
launched on running `dapr init`.

<img src="img/workflow-trace-spans-zipkin.png">

### What happened? 

When you ran the above comands:

1. First the "user" inputs a request to onboard the new Organization TestOrg01 into the concole app.
2. A unique order ID for the workflow is generated (in the above example, `ceb60a94-e598-4b2f-af7f-8ac7e5c683ec`) and the workflow is scheduled.
3. The `NotifyActivity` workflow activity sends a notification saying a request for onboarding TestOrg01 has been received.
4. The `VerifyOnboardingActivity` workflow activity checks that the onboarding data is correct, and confirms it is ready for processing.
5. The `RequestApprovalActivity` workflow activity is triggered due to buisness logic for organisations with more than 50 users and the user is prompted to manually approve the 
purchase before continuing the order.
6. The workflow starts and notifies you of its status.
7. The `ProcessAddUserActivity` workflow activity begins processing user details for onboarding `ceb60a94-e598-4b2f-af7f-8ac7e5c683ec` and confirms if successful.
8. The `UpdateOnboardingActivity` workflow activity updates the status of the onboarding, and confirms the process has completed.
9. The `NotifyActivity` workflow activity sends a notification saying that the onboarding for `ceb60a94-e598-4b2f-af7f-8ac7e5c683ec` has completed.
10. The workflow terminates as completed.

## Observations

- Dapr is code heavy, other workflow orchestration tools might be more suitable for building solutions quickly.
- Dapr is flexible, but development needs to be carefully planned out, otherwise code can get messy and hard to read.

## Redis

- Dapr manual process example

  ```sh
  # Step 1: Run the Dapr sidecar
  dapr run --app-id myapp --dapr-http-port 3500

  # Step 2: Save state
  curl -X POST -H "Content-Type: application/json" -d '[{ "key": "name", "value": "Bruce Wayne"}]' http://localhost:3500/v1.0/state/statestore

  # Step 3: Get state
  curl http://localhost:3500/v1.0/state/statestore/name

  # Step 4: See how the state is stored in Redis
  docker exec -it dapr_redis redis-cli

      ## List the Redis keys to see how Dapr created a key value pair with the app-id you provided to dapr run as the key’s prefix:
      keys *

      ## View the state values by running:
      hgetall "myapp||name"

      ## Exit the Redis CLI with:
      exit

  # Step 5: Delete state
  curl -v -X DELETE -H "Content-Type: application/json" http://localhost:3500/v1.0/state/statestore/name
  ```

- [Clearing Redis Cache: Step-by-Step Guide to Safely Remove Data Stored in Redis](https://copperchips.com/clearing-redis-cache-step-by-step-guide-to-safely-remove-data-stored-in-redis/)

  ```sh
  # select the database_number
  select 0

  # clear the cache
  flushdb

  # clear the cache
  flushall
  ```

## References

- [Github - Organization Service](https://github.com/bytesuk/organization-service/tree/feature/user-service-integration/src)
- [Zipkin - Dapr Query Visualiser Tool](http://localhost:9411/zipkin/)
- [Dapr Github - Quickstart Workflows](https://github.com/dapr/quickstarts/tree/master/workflows/go/sdk)
- [Dapr Docs - Workflow Quickstart](https://docs.dapr.io/getting-started/quickstarts/workflow-quickstart/)
- [Dapr Docs - Workflow Overview](https://docs.dapr.io/developing-applications/building-blocks/workflow/workflow-overview/)
- [Dapr Docs - Go SDK](https://docs.dapr.io/developing-applications/sdks/go/)
- [How to update the Go version](https://gist.github.com/nikhita/432436d570b89cab172dcf2894465753)

## Resources

- [Golangbot - Golang tutorial series](https://golangbot.com/learn-golang-series/)
  - [Golangbot - Variables](https://golangbot.com/variables)
- [How To Use Dates and Times in Go](https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go)
- [Go Docs - Time.IsZero](https://pkg.go.dev/time#Time.IsZero)
- [Embedding structs in structs](https://eli.thegreenplace.net/2020/embedding-in-go-part-1-structs-in-structs)
- [Use Environment Variable in your next Golang Project](https://towardsdatascience.com/use-environment-variable-in-your-next-golang-project-39e17c3aaa66)
