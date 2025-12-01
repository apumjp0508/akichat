package signaling

import (
	"encoding/json"
	"fmt"

	"akichat/backend/internal/realtime"
)

// InboundSignal はクライアント→サーバのシグナリングメッセージ形式（既存実装に合わせる）
type InboundSignal struct {
	Type      string          `json:"type"`
	To        uint            `json:"to"`
	SDP       json.RawMessage `json:"sdp,omitempty"`
	Candidate json.RawMessage `json:"candidate,omitempty"`
}

// DeliveryError は宛先への配送に失敗した場合のエラー（元のtypeや宛先を保持）
type DeliveryError struct {
	OriginalType string
	To           uint
	Err          error
}

func (e *DeliveryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("delivery failed: type=%s to=%d: %v", e.OriginalType, e.To, e.Err)
	}
	return fmt.Sprintf("delivery failed: type=%s to=%d", e.OriginalType, e.To)
}

type Service struct {
	RT realtime.Gateway
}

// Handle は受信した生JSONを解釈し、適切な宛先に中継する
// 既存の payload 形式に合わせて map を構築する
func (s *Service) Handle(senderID uint, raw []byte) error {
	var m InboundSignal
	if err := json.Unmarshal(raw, &m); err != nil || m.Type == "" {
		// シグナリングではない（他のメッセージに委譲）
		return nil
	}

	switch m.Type {
	case "webrtc_offer", "webrtc_answer":
		payload := map[string]interface{}{
			"type": m.Type,
			"from": senderID,
			"sdp":  m.SDP,
		}
		if err := s.RT.SendTo(m.To, payload); err != nil {
			return &DeliveryError{OriginalType: m.Type, To: m.To, Err: err}
		}
		return nil

	case "webrtc_ice":
		payload := map[string]interface{}{
			"type":      m.Type,
			"from":      senderID,
			"candidate": m.Candidate,
		}
		if err := s.RT.SendTo(m.To, payload); err != nil {
			return &DeliveryError{OriginalType: m.Type, To: m.To, Err: err}
		}
		return nil

	default:
		// 未知タイプは無視
		return nil
	}
}



