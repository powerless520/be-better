package queue

import "gopkg.in/Shopify/sarama.v1"

type KafkaConsumer struct {
	ready    chan bool
	handler  func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, consumer *Consumer)
	consumer *Consumer
}

// setup is run at the beginning of a new seesion,before ConsumeClaim
func (consumer *KafkaConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// cleanup is run at the end of a session,once all consumClaim groutines have exited
func (consumer *KafkaConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *KafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	consumer.handler(session, claim, consumer.consumer)
	//for message := range claim.Messages() {
	//	if message.Offset > 0 {
	//		messageJson, _ := json.Marshal(message)
	//		global.GVA_LOG.Debug("KafkaConsumeHandle message: ", string(messageJson)+"\n"+string(message.Value))
	//		consumer.handler(session, claim, consumer.consumer)
	//		session.MarkMessage(message, "")
	//	}
	//}
	return nil
}
