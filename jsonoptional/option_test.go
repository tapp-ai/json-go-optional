package jsonoptional_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tapp-ai/json-go-optional/jsonoptional"
)

type TestRequest struct {
	Value jsonoptional.Option[time.Time] `json:"value,omitempty"`
}

func Test(t *testing.T) {
	t.Run("None", func(t *testing.T) {
		o := jsonoptional.None[time.Time]()
		req := TestRequest{o}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.JSONEq(t, `{}`, string(data))
	})

	t.Run("Some", func(t *testing.T) {
		value := time.Date(2024, 9, 13, 0, 0, 0, 0, time.UTC)
		o := jsonoptional.Some(value)
		req := TestRequest{o}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"value": "2024-09-13T00:00:00Z"}`, string(data))
	})

	t.Run("Null", func(t *testing.T) {
		o := jsonoptional.Null[time.Time]()
		req := TestRequest{o}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"value": null}`, string(data))
	})

	t.Run("NullIf", func(t *testing.T) {
		value := time.Time{}
		o := jsonoptional.NullIf(value, value.IsZero())
		req := TestRequest{o}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"value": null}`, string(data))
	})
}
