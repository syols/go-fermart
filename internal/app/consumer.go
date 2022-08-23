package app

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg/database"
	"github.com/syols/go-devops/internal/pkg/event"
)

func Consume(ctx context.Context, wg *sync.WaitGroup, settings config.Config) error {
	log.Print("CONSUME")
	defer wg.Done()
	sess, err := event.NewSession()
	if err != nil {
		return err
	}

	connection, err := database.NewConnection(settings)
	if err != nil {
		return err
	}

	pollInterval := time.NewTicker(time.Second)
	for {
		select {
		case <-pollInterval.C:
			messages, err := sess.ReceiveMessages()
			if err != nil {
				log.Print(err.Error())
				continue
			}

			for _, msg := range messages {
				url := settings.AccrualAddress + "/api/orders/" + *msg.Body
				resp, err := http.Get(url)
				if err != nil {
					log.Print(err.Error())
					continue
				}

				value := models.Order{}
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Print(err.Error())
					continue
				}

				if err := resp.Body.Close(); err != nil {
					log.Print(err.Error())
				}

				bodyString := string(bodyBytes)
				log.Print(bodyString)

				if err := json.Unmarshal(bodyBytes, &value); err != nil {
					log.Print(err.Error())
					continue
				}

				if value.Status == models.ProcessedOrderStatus {
					if err := sess.DeleteMessage(msg); err != nil {
						log.Print(err.Error())
					}
				}

				if err = value.Update(ctx, connection); err != nil {
					log.Print(err.Error())
					continue
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}
