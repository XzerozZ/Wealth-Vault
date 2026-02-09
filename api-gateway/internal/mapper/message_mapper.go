package mapper

import (
	"encoding/json"
	"wealth-vault/api-gateway/internal/domain"
	pb "wealth-vault/api-gateway/pkg/pb/proto/user"
)

func ToMessageResponseList(pbMsgs []*pb.MessageDetail) []domain.Message {
	var responses []domain.Message

	for _, m := range pbMsgs {
		var metaMap map[string]interface{}
		if m.Metadata != "" {
			_ = json.Unmarshal([]byte(m.Metadata), &metaMap)
		}
		if metaMap == nil {
			metaMap = make(map[string]interface{})
		}

		responses = append(responses, domain.Message{
			ID:          m.Id,
			SenderID:    m.SenderId,
			MsgType:     m.MsgType,
			Content:     m.Content,
			Metadata:    metaMap,
			CreatedAt:   m.CreatedAt.AsTime(),
			SenderName:  m.SenderName,
			SenderImage: m.SenderImage,
			IsMe:        m.IsMe,
		})
	}

	return responses
}
