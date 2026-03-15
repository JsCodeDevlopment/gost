package messaging

import (
	"encoding/json"
	"log"
	"time"

	"gost/src/common/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type WebhookRetryPayload struct {
	URL        string      `json:"url"`
	Secret     string      `json:"secret"`
	Event      string      `json:"event"`
	Data       interface{} `json:"data"`
	RetryCount int         `json:"retry_count"`
}

const (
	WebhookQueue      = "webhook_delivery"
	MaxWebhookRetries = 5
)

func DispatchWebhookWithRetry(url, secret, event string, data interface{}) {
	err := utils.SendWebhook(url, secret, event, data)
	if err != nil {
		log.Printf("Webhook first attempt failed for %s: %v. Queuing for retry...", url, err)

		payload := WebhookRetryPayload{
			URL:        url,
			Secret:     secret,
			Event:      event,
			Data:       data,
			RetryCount: 1,
		}

		PublishMessage("", WebhookQueue, payload)
	}
}

func StartWebhookWorker() {
	err := RegisterConsumer(WebhookQueue, func(d amqp.Delivery) {
		var payload WebhookRetryPayload
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("Failed to unmarshal webhook retry payload: %v", err)
			d.Ack(false)
			return
		}

		log.Printf("Retrying webhook delivery to %s (Attempt %d)", payload.URL, payload.RetryCount)

		err := utils.SendWebhook(payload.URL, payload.Secret, payload.Event, payload.Data)
		if err != nil {
			if payload.RetryCount < MaxWebhookRetries {
				payload.RetryCount++
				time.Sleep(time.Duration(payload.RetryCount) * time.Second)
				PublishMessage("", WebhookQueue, payload)
				log.Printf("Webhook retry %d failed, re-queued", payload.RetryCount)
			} else {
				log.Printf("Webhook failed after maximum retries: %s", payload.URL)
			}
		} else {
			log.Printf("Webhook delivered successfully on retry to %s", payload.URL)
		}

		d.Ack(false)
	})

	if err != nil {
		log.Printf("Failed to start webhook worker: %v", err)
	}
}
