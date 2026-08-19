package service

import (
	"context"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/contest/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type Detail struct {
	ProblemID       uuid.UUID `json:"problem_id"`
	ProblemPosition int       `json:"problem_position"`
	SubmissionCount int       `json:"submission_count"`
	Score           int       `json:"score"`
	LastSubmit      time.Time `json:"last_submit"`
}

type Ranking struct {
	Name    string   `json:"name"`
	Score   float64  `json:"score"`
	Penalty float64  `json:"penalty"`
	Details []Detail `json:"details"`
}

type submissionRow = repository.GetContestSubmissionsForRankingRow

const (
	PenaltyPerSubmission = 10
	ScorePerProblem      = 1
)

func (s *Service) GetRank(
	ctx context.Context,
	contestID uuid.UUID,
) ([]Ranking, error) {
	queries := repository.New(s.pool)

	contestTime, err := queries.GetContestTimeRange(ctx, contestID)
	if err != nil {
		return nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to get contest time range",
			err,
		)
	}

	rows, err := queries.GetContestSubmissionsForRanking(ctx, contestID)
	if err != nil {
		return nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to get contest ranking",
			err,
		)
	}

	result := buildRankings(rows, contestTime.StartTime.Time)

	sort.Slice(result, func(i, j int) bool {
		if result[i].Score != result[j].Score {
			return result[i].Score > result[j].Score
		}

		return result[i].Penalty < result[j].Penalty
	})

	return result, nil
}

func buildRankings(rows []submissionRow, contestStart time.Time) []Ranking {
	result := make([]Ranking, 0)

	for i := 0; i < len(rows); {
		accountID := rows[i].AccountID
		accountName := rows[i].AccountName

		accountEnd := i
		for accountEnd < len(rows) && rows[accountEnd].AccountID == accountID {
			accountEnd++
		}

		details, score, penalty := buildAccountRanking(rows[i:accountEnd], contestStart)

		result = append(result, Ranking{
			Name:    accountName,
			Score:   score,
			Penalty: penalty,
			Details: details,
		})

		i = accountEnd
	}

	return result
}

func buildAccountRanking(rows []submissionRow, contestStart time.Time) ([]Detail, float64, float64) {
	details := make([]Detail, 0)

	var totalScore, totalPenalty float64

	for i := 0; i < len(rows); {
		problemID := rows[i].ProblemID

		problemEnd := i
		for problemEnd < len(rows) && rows[problemEnd].ProblemID == problemID {
			problemEnd++
		}

		detail, penalty := calculateProblemResult(rows[i:problemEnd], contestStart)

		details = append(details, detail)
		totalScore += float64(detail.Score)
		totalPenalty += penalty

		i = problemEnd
	}

	return details, totalScore, totalPenalty
}

func calculateProblemResult(
	submissions []submissionRow,
	contestStart time.Time,
) (Detail, float64) {
	if submissions[0].Language == repository.LanguageHtml {
		return calculateCssProblemResult(submissions)
	}

	return calculateAlgorithmProblemResult(submissions, contestStart)
}

func calculateCssProblemResult(submissions []submissionRow) (Detail, float64) {
	last := submissions[len(submissions)-1]
	submissionCount := len(submissions)

	var best float32
	for _, sub := range submissions {
		if sub.Score != nil && *sub.Score > best {
			best = *sub.Score
		}
	}

	detail := Detail{
		ProblemID:       last.ProblemID,
		ProblemPosition: int(*last.ProblemPosition),
		SubmissionCount: submissionCount,
		Score:           int(math.Round(float64(best))),
		LastSubmit:      last.CreatedAt.Time,
	}

	penalty := float64(submissionCount * PenaltyPerSubmission)

	return detail, penalty
}

func calculateAlgorithmProblemResult(
	submissions []submissionRow,
	contestStart time.Time,
) (Detail, float64) {
	truncated, hasAccepted := truncateAtFirstAccepted(submissions)

	last := truncated[len(truncated)-1]

	submissionCount := len(truncated)

	score := 0
	if hasAccepted {
		score = ScorePerProblem
	}

	lastSubmit := last.CreatedAt.Time

	detail := Detail{
		ProblemID:       last.ProblemID,
		ProblemPosition: int(*last.ProblemPosition),
		SubmissionCount: submissionCount,
		Score:           score,
		LastSubmit:      lastSubmit,
	}

	penalty := 0.0
	if hasAccepted {
		minutes := lastSubmit.Sub(contestStart).Minutes()
		penalty = minutes + float64(submissionCount*PenaltyPerSubmission)
	}

	return detail, penalty
}

func truncateAtFirstAccepted(
	submissions []submissionRow,
) ([]submissionRow, bool) {
	for i, sub := range submissions {
		if sub.Verdict != nil && *sub.Verdict == repository.VerdictAccepted {
			return submissions[:i+1], true
		}
	}

	return submissions, false
}
