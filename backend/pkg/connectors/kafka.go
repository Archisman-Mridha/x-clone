package connectors

import (
	"context"
	"log/slog"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/assert"
	"github.com/Archisman-Mridha/x-clone/backend/pkg/utils"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
	"github.com/twmb/franz-go/plugin/kslog"
)

type (
	KafkaConnector struct {
		client *kgo.Client
	}

	NewKafkaConnectorArgs struct {
		SeedBrokerURLs []string
		Group          string
		Topics         []string
	}
)

func NewKafkaConnector(ctx context.Context, args *NewKafkaConnectorArgs) *KafkaConnector {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(args.SeedBrokerURLs...),

		kgo.ConsumerGroup(args.Group),

		kgo.ConsumeTopics(args.Topics...),
		kgo.AllowAutoTopicCreation(),

		kgo.DisableAutoCommit(),
		kgo.AutoCommitMarks(),

		kgo.WithHooks(
			kslog.New(slog.Default()),
			kotel.WithMeter(kotel.NewMeter()),
			kotel.WithTracer(kotel.NewTracer()),
		),
	)
	assert.AssertErrNil(ctx, err, "Failed connecting to Kafka")

	// Ping Kafka, verifying that a working connection has been established.
	err = client.Ping(ctx)
	assert.AssertErrNil(ctx, err, "Failed connecting to Kafka")

	slog.DebugContext(ctx, "Connected to Kafka")

	return &KafkaConnector{client}
}

func (k *KafkaConnector) GetClient() *kgo.Client {
	return k.client
}

func (k *KafkaConnector) Healthcheck() error {
	err := k.client.Ping(context.Background())
	if err != nil {
		return utils.WrapErrorWithPrefix("Failed pinging Kafka", err)
	}
	return nil
}

func (k *KafkaConnector) Shutdown(ctx context.Context) {
	k.client.Close()
	slog.DebugContext(ctx, "Shut down Kafka client")
}
