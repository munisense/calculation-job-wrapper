package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

type AMQPManager struct {
	channel       *amqp.Channel
	conn          *amqp.Connection
	closeNotifier chan *amqp.Error
	config        *MQConfig
	appId         string
}

func NewAMQPManager(config *MQConfig, appId string) (*AMQPManager, error) {
	man := AMQPManager{}
	man.config = config
	man.appId = appId

	if err := man.initializeChannel(); err != nil {
		return nil, err
	}

	man.closeNotifier = make(chan *amqp.Error)
	man.conn.NotifyClose(man.closeNotifier)
	man.channel.NotifyClose(man.closeNotifier)

	return &man, nil
}

func (man *AMQPManager) initializeChannel() error {
	var err error
	url := fmt.Sprintf("amqps://%s:%s@%s:%d/", man.config.Username, man.config.Password, man.config.Host, man.config.Port)
	man.conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("Could not connect to AMQP: %v", err)
	}

	man.channel, err = man.conn.Channel()

	if err != nil {
		return fmt.Errorf("Could not create AMQP Channel: %v", err)
	}

	man.channel.Qos(5, 0, false)

	return nil
}

func (man *AMQPManager) consume(queue string) (<-chan amqp.Delivery, error) {
	return man.channel.Consume(queue, man.appId, false, true, false, true, nil)
}
