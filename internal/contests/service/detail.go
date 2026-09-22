package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/contests/repository"
	apperr "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ContestProblem struct {
	ID       uuid.UUID `json:"id"`
	Position int       `json:"position"`
	Name     string    `json:"name"`
}

type ContestDetail struct {
	ID       uuid.UUID        `json:"id"`
	Name     string           `json:"name"`
	Start    time.Time        `json:"start"`
	End      time.Time        `json:"end"`
	Problems []ContestProblem `json:"problems"`
}

func (s *Service) GetContestDetail(
	ctx context.Context,
	contestID uuid.UUID,
) (ContestDetail, error) {
	queries := repository.New(s.pool)

	contest, err := queries.GetContest(ctx, contestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ContestDetail{}, apperr.Wrap(
				http.StatusNotFound,
				"contest not found",
				err,
			)
		}

		return ContestDetail{}, apperr.Wrap(
			http.StatusInternalServerError,
			"failed to get contest",
			err,
		)
	}

	now := time.Now()
	if now.Before(contest.Start) {
		return ContestDetail{}, apperr.Wrap(http.StatusBadRequest, "Contest not start yet", nil)
	}

	problems, err := queries.GetProblemsByContest(ctx, contestID)
	if err != nil {
		return ContestDetail{}, apperr.Wrap(
			http.StatusInternalServerError,
			"failed to get contest problems",
			err,
		)
	}

	problemList := make([]ContestProblem, 0, len(problems))
	for _, p := range problems {
		problemList = append(problemList, ContestProblem{
			ID:       p.ID,
			Position: int(*p.Position),
			Name:     p.Name,
		})
	}

	return ContestDetail{
		ID:       contest.ID,
		Name:     contest.Name,
		Start:    contest.Start,
		End:      contest.End,
		Problems: problemList,
	}, nil
}
