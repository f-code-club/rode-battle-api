package service

import (
	"context"
	"net/http"
	"sort"
	"time"
	"unicode"

	"github.com/f-code-club/rode-battle-api/internal/contests/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type Detail struct {
	ProblemID       uuid.UUID `json:"problem_id"`
	ProblemPosition int       `json:"problem_position"`
	SubmissionCount int       `json:"submission_count"`
	Score           float64   `json:"score"`
	LastSubmit      time.Time `json:"last_submit"`
}

type Ranking struct {
	Name    string   `json:"name"`
	Score   float64  `json:"score"`
	Penalty int      `json:"penalty"`
	Details []Detail `json:"details"`
}

type submissionRow = repository.GetContestSubmissionsRow

const (
	PenaltyPerSubmission = 10
	ScorePerProblem      = 1
	PenaltyPerCodeChar   = 1
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

	rows, err := queries.GetContestSubmissions(ctx, contestID)
	if err != nil {
		return nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to get contest ranking",
			err,
		)
	}

	result := buildRankings(rows, contestTime.StartTime)

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

func buildAccountRanking(rows []submissionRow, contestStart time.Time) ([]Detail, float64, int) {
	details := make([]Detail, 0)

	var totalScore float64
	var totalPenalty int

	for i := 0; i < len(rows); {
		problemID := rows[i].ProblemID

		problemEnd := i
		for problemEnd < len(rows) && rows[problemEnd].ProblemID == problemID {
			problemEnd++
		}

		detail, penalty := calculateProblemResult(rows[i:problemEnd], contestStart)

		details = append(details, detail)
		totalScore += detail.Score
		totalPenalty += penalty

		i = problemEnd
	}

	return details, totalScore, totalPenalty
}

func calculateProblemResult(
	submissions []submissionRow,
	contestStart time.Time,
) (Detail, int) {
	if submissions[0].Language == repository.LanguageHtml {
		return calculateCssProblemResult(submissions, contestStart)
	}

	return calculateAlgorithmProblemResult(submissions, contestStart)
}

func calculateCssProblemResult(submissions []submissionRow, contestStart time.Time) (Detail, int) {
	last := submissions[len(submissions)-1]
	submissionCount := len(submissions)

	var best float64
	var bestCode string
	for _, sub := range submissions {
		if sub.Score != nil && float64(*sub.Score) > best {
			best = float64(*sub.Score)
			bestCode = sub.Code
		}
	}

	detail := Detail{
		ProblemID:       last.ProblemID,
		ProblemPosition: int(*last.ProblemPosition),
		SubmissionCount: submissionCount,
		Score:           best,
		LastSubmit:      last.CreatedAt,
	}

	codeLength := effectiveCSSLength(bestCode)

	minutes := last.CreatedAt.Sub(contestStart).Minutes()
	penalty := int(minutes) + submissionCount*PenaltyPerSubmission + codeLength*PenaltyPerCodeChar

	return detail, penalty
}

func calculateAlgorithmProblemResult(
	submissions []submissionRow,
	contestStart time.Time,
) (Detail, int) {
	truncated, hasAccepted := truncateAtFirstAccepted(submissions)

	last := truncated[len(truncated)-1]

	submissionCount := len(truncated)

	var score float64
	if hasAccepted {
		score = ScorePerProblem
	}

	lastSubmit := last.CreatedAt

	detail := Detail{
		ProblemID:       last.ProblemID,
		ProblemPosition: int(*last.ProblemPosition),
		SubmissionCount: submissionCount,
		Score:           score,
		LastSubmit:      lastSubmit,
	}

	penalty := 0
	if hasAccepted {
		minutes := lastSubmit.Sub(contestStart).Minutes()
		penalty = int(minutes) + submissionCount*PenaltyPerSubmission
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

func effectiveCSSLength(code string) int {
	count := 0
	for _, r := range code {
		if unicode.IsSpace(r) {
			continue
		}
		count++
	}
	return count
}
