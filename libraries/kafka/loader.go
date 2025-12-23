package kafka

import (
	"github.com/semanggilab/webcore-go/app/config"
	"github.com/semanggilab/webcore-go/app/loader"
)

type KafkaConsumerLoader struct {
	name string
}

func (a *KafkaConsumerLoader) SetName(name string) {
	a.name = name
}

func (a *KafkaConsumerLoader) Name() string {
	return a.name
}

func (l *KafkaConsumerLoader) Init(args ...any) (loader.Library, error) {
	config := args[0].(config.KafkaConfig)
	receiver := args[1].(KafkaReceiver)

	kc, err := NewKafkaConsumer(&config, receiver)
	if err != nil {
		return nil, err
	}

	err = kc.Install(args...)
	if err != nil {
		return nil, err
	}

	kc.Connect()
	return kc, nil
}

type KafkaProducerLoader struct {
	name string
}

func (a *KafkaProducerLoader) SetName(name string) {
	a.name = name
}

func (a *KafkaProducerLoader) Name() string {
	return a.name
}

func (l *KafkaProducerLoader) Init(args ...any) (loader.Library, error) {
	config := args[0].(config.KafkaConfig)

	kc, err := NewKafkaProducer(&config)
	if err != nil {
		return nil, err
	}

	err = kc.Install(args...)
	if err != nil {
		return nil, err
	}

	kc.Connect()
	return kc, nil
}
