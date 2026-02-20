package worker

import (
	"context"
	"encoding/json"
)

func wrap[T any](
	handler func(context.Context, T) error,
) func(context.Context, []byte) error {

	return func(ctx context.Context, data []byte) error {
		var evt T

		if err := json.Unmarshal(data, &evt); err != nil {
			return err
		}

		return handler(ctx, evt)
	}
}
