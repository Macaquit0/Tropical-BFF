package hub

import (
	"context"
	"encoding/json"

	"github.com/backend/bff-cognito/pkg/internalaws"
	"github.com/backend/bff-cognito/pkg/logger"
)

type Publisher struct {
	log       *logger.Logger
	publisher internalaws.AwsSqsPublisher
}

func NewPublisher(log *logger.Logger, awsCfg internalaws.Config) (*Publisher, error) {
	publisher, err := internalaws.NewPublisher(log, awsCfg)
	if err != nil {
		return nil, err
	}
	return &Publisher{
		log,
		publisher,
	}, nil
}

type NewProductParams struct {
	TaxId   string `json:"tax_id"`
	Product string `json:"product"`
}

func (n *Publisher) PublishNewProduct(ctx context.Context, product, taxId string) error {

	data := NewProductParams{
		TaxId:   taxId,
		Product: product,
	}

	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if err := n.publisher.Publish(ctx, internalaws.PublishMessage{
		Type: "handle_notification",
		Message: internalaws.PublisherInput{
			Version: 1,
			Data:    b,
		},
	}); err != nil {
		return err
	}

	return nil
}
