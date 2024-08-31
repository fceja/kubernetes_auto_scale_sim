package lib

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// Process message.
func ProcessMessage(msg *sarama.ConsumerMessage) {
	zap.L().Debug("Processing claim message.",
		zap.String("sleeping", fmt.Sprintf("%v", 5*time.Second)),
		zap.String("partition", fmt.Sprintf("%+v", msg.Partition)),
		zap.String("offset", fmt.Sprintf("%+v", msg.Offset)),
		zap.String("jsonEncoded", string(msg.Value)),
	)

	// sleep - simulate work
	time.Sleep(3 * time.Second)
}
