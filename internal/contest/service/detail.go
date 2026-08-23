package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	contestrepo "github.com/f-code-club/rode-battle-api/internal/contest/repository"
	problemrepo "github.com/f-code-club/rode-battle-api/internal/problems/repository"
	sharederrors "github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Problem struct {
	ID       uuid.UUID `json:"id"`
	Position int       `json:"position"`
	Name     string    `json:"name"`
}

type ContestDetail struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Problems []Problem `json:"problems"`
}

func (s *Service) GetContestDetail(
	ctx context.Context,
	contestID uuid.UUID,
) (ContestDetail, error) {
	contestQueries := contestrepo.New(s.pool)
	problemQueries := problemrepo.New(s.pool)

	contest, err := contestQueries.GetContest(ctx, contestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ContestDetail{}, sharederrors.Wrap(
				http.StatusNotFound,
				"contest not found",
				err,
			)
		}

		return ContestDetail{}, sharederrors.Wrap(
			http.StatusInternalServerError,
			"failed to get contest",
			err,
		)
	}

	problems, err := problemQueries.GetProblemsByContest(ctx, contestID)
	if err != nil {
		return ContestDetail{}, sharederrors.Wrap(
			http.StatusInternalServerError,
			"failed to get contest problems",
			err,
		)
	}

	problemList := make([]Problem, 0, len(problems))
	for _, p := range problems {
		problemList = append(problemList, Problem{
			ID:       p.ID,
			Position: int(*p.Position),
			Name:     p.Name,
		})
	}

	return ContestDetail{
		ID:       contest.ID,
		Name:     contest.Name,
		Start:    contest.Start.Time,
		End:      contest.End.Time,
		Problems: problemList,
	}, nil
}
