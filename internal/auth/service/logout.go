package service

import "context"

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.refreshTokenSvc.ParseToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}
