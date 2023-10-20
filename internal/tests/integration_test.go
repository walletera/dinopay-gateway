package tests

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    "github.com/walletera/message-processor/pkg/rabbitmq"
    "net/http"
    "net/url"
    "testing"
    "time"

    msClient "github.com/walletera/mockserver-go-client/pkg/client"
    "go.uber.org/zap"
)

const (
    mockserverPort         = "2090"
    testTimeout            = 30 * time.Second
    testTerminationTimeout = 10 * time.Second
)

type MockServerExpectation struct {
    ExpectationID string `json:"id"`
}

type integrationTest struct {
    t                *testing.T
    mockserverClient *msClient.Client
    expectationIds   []string
    logger           *zap.Logger
}

func newIntegrationTest(t *testing.T, testName string, expectations ...[]byte) *integrationTest {

    logger, _ := zap.NewDevelopment()

    logger.Info("starting test containers")

    initializationCtx, cancelInitializationCtx := context.WithTimeout(context.Background(), testTimeout)
    defer cancelInitializationCtx()

    startRabbitMQContainer(initializationCtx, t, logger)
    startMockserverContainer(initializationCtx, t, logger)

    mockserverUrl, err := url.Parse(fmt.Sprintf("http://localhost:%s", mockserverPort))
    require.NoError(t, err)

    mockserverClient := msClient.NewClient(mockserverUrl, http.DefaultClient)

    logger.Info("creating expectations")

    var expectationsIds = make([]string, len(expectations))
    for i, expectation := range expectations {
        var unmarshalledExpectation MockServerExpectation
        err := json.Unmarshal(expectation, &unmarshalledExpectation)
        require.NoError(t, err, "unmarshalling expectation")
        expectationsIds[i] = unmarshalledExpectation.ExpectationID
        err = mockserverClient.CreateExpectation(initializationCtx, expectation)
        require.NoError(t, err, "creating mockserver expectations")
    }

    return &integrationTest{
        t:                t,
        mockserverClient: mockserverClient,
        expectationIds:   expectationsIds,
        logger:           logger,
    }
}

func (it *integrationTest) Run(test func(ctx context.Context, t *testing.T)) {
    testCtx, _ := context.WithTimeout(context.Background(), testTimeout)

    testFinished := make(chan any)
    go func() {
        test(testCtx, it.t)
        close(testFinished)
    }()

    select {
    case <-testFinished:
    case <-testCtx.Done():
        it.t.Fatal("timeout waiting for test to finish")
    }

    it.logger.Info("verifying expectations")
    for _, id := range it.expectationIds {
        err := it.mockserverClient.VerifyRequest(testCtx, msClient.VerifyRequestBody{ExpectationId: msClient.ExpectationId{
            Id: id,
        }})
        require.NoError(it.t, err, "verification for expectation id %s", id)
    }
}

func startRabbitMQContainer(ctx context.Context, t *testing.T, logger *zap.Logger) {
    req := testcontainers.ContainerRequest{
        Image: "rabbitmq:3.8.0-management",
        Name:  "rabbitmq",
        User:  "rabbitmq",
        ExposedPorts: []string{
            fmt.Sprintf("%d:%d", rabbitmq.DefaultPort, rabbitmq.DefaultPort),
            fmt.Sprintf("%d:%d", rabbitmq.ManagementUIPort, rabbitmq.ManagementUIPort),
        },
        WaitingFor: wait.NewExecStrategy([]string{"rabbitmqadmin", "list", "queues"}).WithStartupTimeout(20 * time.Second),
    }
    rabbitmqC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err, "error creating rabbitmq container")
    t.Cleanup(func() {
        terminationCtx, terminationCtxCancel := context.WithTimeout(context.Background(), testTerminationTimeout)
        defer terminationCtxCancel()
        terminationErr := rabbitmqC.Terminate(terminationCtx)
        if terminationErr != nil {
            logger.Error("failed terminating rabbitmq container", zap.Error(terminationErr))
        }
    })
}

func startMockserverContainer(ctx context.Context, t *testing.T, logger *zap.Logger) {
    req := testcontainers.ContainerRequest{
        Image: "mockserver/mockserver",
        Name:  "mockserver",
        Env: map[string]string{
            "MOCKSERVER_SERVER_PORT": mockserverPort,
            "MOCKSERVER_LOG_LEVEL":   "DEBUG",
        },
        ExposedPorts: []string{fmt.Sprintf("%s:%s", mockserverPort, mockserverPort)},
        WaitingFor:   wait.ForHTTP("/mockserver/status").WithMethod(http.MethodPut).WithPort(mockserverPort),
    }
    mockserverC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    require.NoError(t, err, "error creating mockserver container")
    t.Cleanup(func() {
        terminationCtx, terminationCtxCancel := context.WithTimeout(context.Background(), testTerminationTimeout)
        defer terminationCtxCancel()
        terminationErr := mockserverC.Terminate(terminationCtx)
        if terminationErr != nil {
            logger.Error("failed terminating mockserver container", zap.Error(err))
        }
    })
}
