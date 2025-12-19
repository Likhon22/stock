package bootstrap

import (
	"context"
	"price_generator/kafka"
)

type Application struct {
	producer *kafka.Producer
}

func InitializeApp(ctx context.Context) (*Application, error) {
	producer := kafka.NewProducer()

}

func (app *Application) Run() {

}
