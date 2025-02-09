package event

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ConsumerConfig struct {
	URL                string
	ReconnectDelay     time.Duration
	GetLatestTimestamp func() string
	MaxRetries         int
}

func ConsumeEvents[T any](
	ctx context.Context, config ConsumerConfig, eventChan chan<- T, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(eventChan)

	retries := 0
	isReconnecting := false

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context canceled, stopping consumer")
			return
		default:
			url := config.URL

			if isReconnecting {
				timestamp := config.GetLatestTimestamp()
				if timestamp != "" {
					log.Printf("Reconnecting: Using timestamp %s", timestamp)

					if strings.Contains(url, "?") {
						url += "&since=" + timestamp
					} else {
						url += "?since=" + timestamp
					}
				}
			}

			err := consumeEventStream(ctx, url, eventChan)

			if err != nil {
				retries++
				if config.MaxRetries > 0 && retries >= config.MaxRetries {
					log.Printf("Max retries (%d) reached, stopping consumer", config.MaxRetries)
					return
				}

				log.Printf("Stream ended or error occurred: %v", err)
				log.Printf("Attempting reconnection #%d in %v...", retries, config.ReconnectDelay)

				isReconnecting = true

				select {
				case <-ctx.Done():
					return
				case <-time.After(config.ReconnectDelay):
					continue
				}
			}

			retries = 0
			isReconnecting = false
		}
	}
}

func consumeEventStream[T any](ctx context.Context, url string, eventChan chan<- T) error {
	log.Println("Connecting to event stream:", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReaderSize(resp.Body, 16*1024)
	buffer := strings.Builder{}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := reader.ReadString('\n')

			if line != "" {
				buffer.WriteString(line)

				if strings.HasSuffix(line, "\n") {
					event := buffer.String()
					buffer.Reset()

					if err := handleEvent(event, eventChan); err != nil {
						log.Printf("Error processing event: %v", err)
					}
				}
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					return fmt.Errorf("stream ended")
				}

				log.Printf("Error reading: %v", err)

				if buffer.Len() > 0 {
					if err := handleEvent(buffer.String(), eventChan); err != nil {
						log.Printf("Error processing remaining buffer: %v", err)
					}

					buffer.Reset()
				}

				continue
			}
		}
	}
}

func handleEvent[T any](eventStr string, eventChan chan<- T) error {
	eventStr = strings.TrimSpace(eventStr)
	if eventStr == "" {
		return nil
	}

	lines := strings.Split(eventStr, "\n")

	var dataLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			dataLines = append(dataLines, data)
		}
	}

	if len(dataLines) == 0 {
		return nil
	}

	completeData := strings.Join(dataLines, "")

	var event T
	if err := json.Unmarshal([]byte(completeData), &event); err != nil {
		return fmt.Errorf("error unmarshaling event: %w", err)
	}

	eventChan <- event

	return nil
}
