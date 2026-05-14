package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/stepup-ai/ai-service/internal/entity"
)

type AIRepository interface {
	SaveAdmissionAnalysis(ctx context.Context, analysis *entity.AdmissionAnalysis) error
	GetAdmissionAnalysisByID(ctx context.Context, analysisID string) (*entity.AdmissionAnalysis, error)
	GetAdmissionAnalysisHistoryByUserID(ctx context.Context, userID string) ([]*entity.AdmissionAnalysis, error)
	SaveRoadmap(ctx context.Context, roadmap *entity.Roadmap, steps []*entity.RoadmapStep) error
	GetRoadmapByUserID(ctx context.Context, userID string) (*entity.Roadmap, []*entity.RoadmapStep, error)
	SaveEssayReview(ctx context.Context, review *entity.EssayReview) error
	GetEssayReviewsByUserID(ctx context.Context, userID string) ([]*entity.EssayReview, error)
}

type PostgresAIRepository struct {
	db *sql.DB
}

func NewPostgresAIRepository(db *sql.DB) AIRepository {
	return &PostgresAIRepository{db: db}
}

func (r *PostgresAIRepository) SaveAdmissionAnalysis(ctx context.Context, analysis *entity.AdmissionAnalysis) error {
	query := `
    INSERT INTO admission_analyses (id, user_id, university_id, admission_chance_percentage, weak_areas, improvement_suggestions, gap_analysis, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  `
	_, err := r.db.ExecContext(ctx, query,
		analysis.ID, analysis.UserID, analysis.UniversityID,
		analysis.AdmissionChancePercentage,
		pq.Array(analysis.WeakAreas),
		pq.Array(analysis.ImprovementSuggestions),
		analysis.GapAnalysis, analysis.CreatedAt,
	)
	return err
}

func (r *PostgresAIRepository) GetAdmissionAnalysisByID(ctx context.Context, analysisID string) (*entity.AdmissionAnalysis, error) {
	query := `
    SELECT id, user_id, university_id, admission_chance_percentage, weak_areas, improvement_suggestions, gap_analysis, created_at
    FROM admission_analyses WHERE id = $1
  `
	analysis := &entity.AdmissionAnalysis{}
	err := r.db.QueryRowContext(ctx, query, analysisID).Scan(
		&analysis.ID, &analysis.UserID, &analysis.UniversityID,
		&analysis.AdmissionChancePercentage,
		pq.Array(&analysis.WeakAreas),
		pq.Array(&analysis.ImprovementSuggestions),
		&analysis.GapAnalysis, &analysis.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

func (r *PostgresAIRepository) GetAdmissionAnalysisHistoryByUserID(ctx context.Context, userID string) ([]*entity.AdmissionAnalysis, error) {
	query := `
    SELECT id, user_id, university_id, admission_chance_percentage, weak_areas, improvement_suggestions, gap_analysis, created_at
    FROM admission_analyses WHERE user_id = $1 ORDER BY created_at DESC
  `
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*entity.AdmissionAnalysis
	for rows.Next() {
		analysis := &entity.AdmissionAnalysis{}
		if err := rows.Scan(
			&analysis.ID, &analysis.UserID, &analysis.UniversityID,
			&analysis.AdmissionChancePercentage,
			pq.Array(&analysis.WeakAreas),
			pq.Array(&analysis.ImprovementSuggestions),
			&analysis.GapAnalysis, &analysis.CreatedAt,
		); err != nil {
			return nil, err
		}
		analyses = append(analyses, analysis)
	}
	return analyses, nil
}

func (r *PostgresAIRepository) SaveRoadmap(ctx context.Context, roadmap *entity.Roadmap, steps []*entity.RoadmapStep) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	roadmapQuery := `INSERT INTO roadmaps (id, user_id, created_at) VALUES ($1, $2, $3)`
	_, err = tx.ExecContext(ctx, roadmapQuery, roadmap.ID, roadmap.UserID, roadmap.CreatedAt)
	if err != nil {
		return err
	}
	for _, step := range steps {
		stepQuery := `
      INSERT INTO roadmap_steps (id, roadmap_id, title, description, month_number, category)
      VALUES ($1, $2, $3, $4, $5, $6)
    `
		_, err = tx.ExecContext(ctx, stepQuery,
			uuid.New().String(), roadmap.ID,
			step.Title, step.Description, step.MonthNumber, step.Category,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresAIRepository) GetRoadmapByUserID(ctx context.Context, userID string) (*entity.Roadmap, []*entity.RoadmapStep, error) {
	roadmapQuery := `SELECT id, user_id, created_at FROM roadmaps WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	roadmap := &entity.Roadmap{}
	err := r.db.QueryRowContext(ctx, roadmapQuery, userID).Scan(
		&roadmap.ID, &roadmap.UserID, &roadmap.CreatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	stepsQuery := `SELECT id, roadmap_id, title, description, month_number, category FROM roadmap_steps WHERE roadmap_id = $1`
	rows, err := r.db.QueryContext(ctx, stepsQuery, roadmap.ID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var steps []*entity.RoadmapStep
	for rows.Next() {
		step := &entity.RoadmapStep{}
		if err := rows.Scan(
			&step.ID, &step.RoadmapID, &step.Title,
			&step.Description, &step.MonthNumber, &step.Category,
		); err != nil {
			return nil, nil, err
		}
		steps = append(steps, step)
	}
	return roadmap, steps, nil
}

func (r *PostgresAIRepository) SaveEssayReview(ctx context.Context, review *entity.EssayReview) error {
	query := `
    INSERT INTO essay_reviews (id, user_id, essay_text, university_name, program_name, grammar_score, coherence_score, uniqueness_score, relevance_score, improvement_suggestions, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
  `
	_, err := r.db.ExecContext(ctx, query,
		review.ID, review.UserID, review.EssayText,
		review.UniversityName, review.ProgramName,
		review.GrammarScore, review.CoherenceScore,
		review.UniquenessScore, review.RelevanceScore,
		pq.Array(review.ImprovementSuggestions),
		time.Now(),
	)
	return err
}

func (r *PostgresAIRepository) GetEssayReviewsByUserID(ctx context.Context, userID string) ([]*entity.EssayReview, error) {
	query := `
    SELECT id, user_id, essay_text, university_name, program_name, grammar_score, coherence_score, uniqueness_score, relevance_score, improvement_suggestions, created_at
    FROM essay_reviews WHERE user_id = $1 ORDER BY created_at DESC
  `
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*entity.EssayReview
	for rows.Next() {
		review := &entity.EssayReview{}
		if err := rows.Scan(
			&review.ID, &review.UserID, &review.EssayText,
			&review.UniversityName, &review.ProgramName,
			&review.GrammarScore, &review.CoherenceScore,
			&review.UniquenessScore, &review.RelevanceScore,
			pq.Array(&review.ImprovementSuggestions),
			&review.CreatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
