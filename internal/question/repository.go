package question

import (
	"github.com/ShopOnGO/review-service/pkg/db"
)

type QuestionRepository struct {
	Db *db.Db
}

func NewQuestionRepository(db *db.Db) *QuestionRepository {
	return &QuestionRepository{
		Db: db,
	}
}

func (r *QuestionRepository) CreateQuestion(question *Question) error {
	return r.Db.Create(question).Error
}

func (r *QuestionRepository) GetQuestionsByProductVariantID(productVariantID uint) ([]Question, error) {
	var questions []Question
	err := r.Db.Where("product_variant_id = ?", productVariantID).Find(&questions).Error
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (r *QuestionRepository) GetQuestionByID(id uint) (*Question, error) {
	var question Question
	err := r.Db.First(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *QuestionRepository) UpdateQuestion(question *Question) error {
	return r.Db.Save(question).Error
}

func (r *QuestionRepository) UpdateAnswer(questionID uint, answer string) error {
	return r.Db.Model(&Question{}).Where("id = ?", questionID).Update("answer_text", answer).Error
}

func (r *QuestionRepository) DeleteQuestion(question *Question) error {
	return r.Db.Delete(question).Error
}

func (r *QuestionRepository) DeleteQuestionByID(id uint) error {
	return r.Db.Delete(&Question{}, id).Error
}

func (r *QuestionRepository) GetQuestionsByProductVariantIDPaginated(productVariantID uint, limit, offset int) ([]*Question, error) {
    var questions []*Question
    result := r.Db.
        Where("product_variant_id = ?", productVariantID).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&questions)

    return questions, result.Error
}