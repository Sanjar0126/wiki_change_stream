package event

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

func ConsumeEvents[T any](ctx context.Context, url string, eventChan chan<- T, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(eventChan)

	log.Println(url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled, stopping consumer")
			return
		default:
			line, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error reading: %v", err)
				continue
			}

			if strings.TrimSpace(line) == "" {
				continue
			}

			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")

				var event T

				if err := json.Unmarshal([]byte(data), &event); err != nil {
					log.Printf("Error unmarshaling event: %v", err)
					continue
				}
				eventChan <- event
			}
		}
	}
}
