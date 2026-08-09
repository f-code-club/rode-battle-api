package service

import "context"

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	id, err := s.refreshTokenSvc.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	return s.accessTokenSvc.GenerateToken(id)
}
