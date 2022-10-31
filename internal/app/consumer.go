package app

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/pkg"
)

func Consume(ctx context.Context, settings config.Config, errs chan error) {
	sess, err := pkg.NewSession()
	if err != nil {
		errs <- err
		close(errs)
		return
	}

	conn := pkg.NewDatabaseUrlConnection(settings)
	db, err := pkg.NewDatabase(conn)
	if err != nil {
		errs <- err
		close(errs)
		return
	}

	pollInterval := time.NewTicker(time.Second)
	for {
		select {
		case <-pollInterval.C:
			messages, err := sess.ReceiveMessages()
			if err != nil {
				errs <- err
				continue
			}

			for _, msg := range messages {
				url := settings.AccrualAddress + "/api/orders/" + *msg.Body
				resp, err := http.Get(url)
				if err != nil {
					errs <- err
					continue
				}

				value := models.Order{}
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					errs <- err
					continue
				}

				if err := resp.Body.Close(); err != nil {
					errs <- err
				}

				if err := json.Unmarshal(bodyBytes, &value); err != nil {
					errs <- err
					continue
				}

				if value.Status == models.ProcessedOrderStatus {
					if err := sess.DeleteMessage(msg); err != nil {
						errs <- err
					}
				}

				if err = value.Update(ctx, db); err != nil {
					errs <- err
				}
			}
		case <-ctx.Done():
			close(errs)
			return
		}
	}
}
