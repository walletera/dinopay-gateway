package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/walletera/eventskit/eventstoredb"
	"github.com/walletera/eventskit/rabbitmq"
)

const (
	mockserverPort               = "2090"
	eventStoreDBPort             = "2113"
	containersStartTimeout       = 30 * time.Second
	containersTerminationTimeout = 10 * time.Second
)

func TestMain(m *testing.M) {
	ctx, _ := context.WithTimeout(context.Background(), containersStartTimeout)

	terminateEventSToreDBContainer, err := startEventStoreDBContainer(ctx)
	if err != nil {
		panic("error starting esdb container: " + err.Error())
	}

	err = eventstoredb.EnableByCategoryProjection(ctx, eventStoreDBUrl)
	if err != nil {
		panic("failed enabling esdb by category projection: " + err.Error())
	}

	err = eventstoredb.SetESDBByCategoryProjectionSeparator(ctx, eventStoreDBUrl)
	if err != nil {
		panic("failed setting esdb by category projection separator: " + err.Error())
	}

	terminateRabbitMQContainer, err := startRabbitMQContainer(ctx)
	if err != nil {
		panic("error starting rabbitmq container: " + err.Error())
	}

	terminateMockserverContainer, err := startMockserverContainer(ctx)
	if err != nil {
		panic("error starting mockserver container:" + err.Error())
	}

	status := m.Run()

	err = terminateEventSToreDBContainer()
	if err != nil {
		panic("error terminating esdb container: " + err.Error())
	}

	err = terminateRabbitMQContainer()
	if err != nil {
		panic("error terminating rabbitmq container: " + err.Error())
	}

	err = terminateMockserverContainer()
	if err != nil {
		panic("error terminating mockserver container: " + err.Error())
	}

	os.Exit(status)
}

func startEventStoreDBContainer(ctx context.Context) (func() error, error) {
	req := testcontainers.ContainerRequest{
		Image: "eventstore/eventstore:21.10.7-buster-slim",
		Name:  "esdb-node",
		Cmd:   []string{"--insecure", "--run-projections=All"},
		ExposedPorts: []string{
			fmt.Sprintf("%s:%s", eventStoreDBPort, eventStoreDBPort),
		},
		WaitingFor: wait.
			ForHTTP("/health/live").
			WithPort("2113/tcp").
			WithStartupTimeout(10 * time.Second).
			WithStatusCodeMatcher(func(status int) bool {
				return status == http.StatusNoContent
			}),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{NewContainerLogConsumer("esdb")},
		},
	}
	esdbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating rabbitmq container: %w", err)
	}

	return func() error {
		terminationCtx, terminationCtxCancel := context.WithTimeout(context.Background(), containersTerminationTimeout)
		defer terminationCtxCancel()
		terminationErr := esdbContainer.Terminate(terminationCtx)
		if terminationErr != nil {
			fmt.Errorf("failed terminating esdb container: %w", err)
		}
		return nil
	}, nil
}

func startRabbitMQContainer(ctx context.Context) (func() error, error) {
	req := testcontainers.ContainerRequest{
		Image: "rabbitmq:3.8.0-management",
		Name:  "rabbitmq",
		User:  "rabbitmq",
		ExposedPorts: []string{
			fmt.Sprintf("%d:%d", rabbitmq.DefaultPort, rabbitmq.DefaultPort),
			fmt.Sprintf("%d:%d", rabbitmq.ManagementUIPort, rabbitmq.ManagementUIPort),
		},
		WaitingFor: wait.NewExecStrategy([]string{"rabbitmqadmin", "list", "queues"}).WithStartupTimeout(20 * time.Second),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{NewContainerLogConsumer("rabbitmq")},
		},
	}
	rabbitmqC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating rabbitmq container: %w", err)
	}

	return func() error {
		terminationCtx, terminationCtxCancel := context.WithTimeout(context.Background(), containersTerminationTimeout)
		defer terminationCtxCancel()
		terminationErr := rabbitmqC.Terminate(terminationCtx)
		if terminationErr != nil {
			fmt.Errorf("failed terminating rabbitmq container: %w", err)
		}
		return nil
	}, nil
}

func startMockserverContainer(ctx context.Context) (func() error, error) {
	req := testcontainers.ContainerRequest{
		Image: "mockserver/mockserver",
		Name:  "mockserver",
		Env: map[string]string{
			"MOCKSERVER_SERVER_PORT": mockserverPort,
			"MOCKSERVER_LOG_LEVEL":   "DEBUG",
		},
		ExposedPorts: []string{fmt.Sprintf("%s:%s", mockserverPort, mockserverPort)},
		WaitingFor:   wait.ForHTTP("/mockserver/status").WithMethod(http.MethodPut).WithPort(mockserverPort),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{NewContainerLogConsumer("mockserver")},
		},
	}
	mockserverC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating mockserver container: %w", err)
	}

	return func() error {
		terminationCtx, terminationCtxCancel := context.WithTimeout(context.Background(), containersTerminationTimeout)
		defer terminationCtxCancel()
		terminationErr := mockserverC.Terminate(terminationCtx)
		if terminationErr != nil {
			fmt.Errorf("failed terminating mockserver container: %w", err)
		}
		return nil
	}, nil
}
