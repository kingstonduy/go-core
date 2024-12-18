package util

import (
	"context"
	"encoding/json"

	"github.com/kingstonduy/go-core/logger"
)

func MakeStringLogs(ctx context.Context, headers map[string][]string, body interface{}, httpStatus string) string {
	object := make(map[string]interface{})

	// Add headers
	if headers != nil {
		object["headers"] = headers
	} else {
		object["headers"] = "empty"
	}

	// Add HTTP status if not empty
	if httpStatus != "" {
		object["http_status"] = httpStatus
	}

	// Add body
	if body != nil {
		if _, ok := body.([]byte); !ok { // Check if body is of type []byte (Buffer)
			object["body"] = body
		} else {
			var jsonBody interface{}
			if err := json.Unmarshal(body.([]byte), &jsonBody); err == nil {
				object["body"] = jsonBody
			} else {
				object["body"] = string(body.([]byte)) // Convert []byte to string
			}
		}
	} else {
		object["body"] = "empty"
	}

	// Convert to JSON string and remove newlines
	jsonString, err := json.Marshal(object)
	if err != nil {
		logger.Errorf(ctx, "Error marshaling to JSON: %v\n", err)
		return ""
	}

	return string(jsonString)
}
