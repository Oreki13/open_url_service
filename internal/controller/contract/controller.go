package contract

import (
	"context"
	"open_url_service/internal/appctx"
	"open_url_service/pkg/pubsubx"
)

type PubSubMessageController interface {
	Serve(ctx context.Context, message *pubsubx.Message)
}

// RabitMq
//type MessageController interface {
//	Serve(data amqp091.Delivery) error
//}

type Controller interface {
	Serve(xCtx appctx.Data) appctx.Response
}
