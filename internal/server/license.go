package server

import "github.com/trysourcetool/onprem-portal/internal/core"

type licenseResponse struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Key    string `json:"key"`
}

func (s *Server) licenseFromModel(l *core.License) *licenseResponse {
	if l == nil {
		return nil
	}

	key, err := s.encryptor.Decrypt(l.KeyCiphertext, l.KeyNonce)
	if err != nil {
		return nil
	}

	return &licenseResponse{
		ID:     l.ID.String(),
		UserID: l.UserID.String(),
		Key:    string(key),
	}
}
