package coreservices

import (
	"context"
	"database/sql"

	"github.com/pedrocmart/leaderboard-service/models"
)

func NewStoreService(core *models.Core, db *sql.DB) models.StoreService {
	storeService := BasicStoreService{
		core: core,
	}
	core.StoreService = &storeService
	core.DB = db
	return &storeService
}

type BasicStoreService struct {
	core *models.Core
}

func (b *BasicStoreService) CreateUser(ctx context.Context, id int, total int) error {
	tx, err := b.core.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO users (id, score) VALUES ($1, $2)`,
		id, total)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (b *BasicStoreService) UpdateRelativeUserScore(ctx context.Context, id int, score int) error {
	tx, err := b.core.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE users 
		SET score = score + $1
		WHERE id = $2`, score, id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func (b *BasicStoreService) UpdateAbsoluteUserScore(ctx context.Context, id int, score int) error {
	tx, err := b.core.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE users 
		SET score = $1
		WHERE id = $2`, score, id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return err
}

func (b *BasicStoreService) DoesUserExist(ctx context.Context, id int) (bool, error) {
	if err := b.core.DB.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id).Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func (b *BasicStoreService) GetUsers(ctx context.Context, top int) ([]models.Ranking, error) {
	rows, err := b.core.DB.QueryContext(ctx, "SELECT id, score FROM users ORDER BY score DESC LIMIT $1", top)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	position := 1
	ranking := make([]models.Ranking, 0)

	for rows.Next() {
		var id int
		var score int
		err = rows.Scan(&id, &score)
		if err != nil {
			return nil, err
		}
		ranking = append(ranking, models.Ranking{
			Position: position,
			UserID:   id,
			Score:    score,
		})
		position++
	}

	return ranking, nil
}

func (b *BasicStoreService) GetUsersBetween(ctx context.Context, pos, around int) ([]models.Ranking, error) {
	/*
		offset is zero-based, and for this reason -1 is being subtracted
		eg:
		pos sent was 100
		offset:= 100 - 3 - 1
		offset = 96
		this way I can start from position 97th (zero-based)
	*/
	offset := pos - around - 1
	if offset < 0 {
		offset = 0
	}
	/*
		since we need the upper, lower and the exactly positions, I' duplicating the parameter and adding 1
		eg:
		positionAround sent was 3
		positionAround = 3 + 3 + 1
		and this will limit the 7 close positions to the position that we want
	*/
	positionAround := around + around + 1
	//if offset is 0, means that we gonna get the top 1 from users; this way we need to get only users below his position, including she/he
	if offset == 0 {
		positionAround = around + 1
	}

	rows, err := b.core.DB.QueryContext(ctx, "SELECT id, score FROM users ORDER BY score DESC LIMIT $1 OFFSET $2", positionAround, offset)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	position := offset + 1 //I'm doing +1 here because offset is zero-based
	ranking := make([]models.Ranking, 0)

	for rows.Next() {
		var id int
		var score int
		err = rows.Scan(&id, &score)
		if err != nil {
			return nil, err
		}
		ranking = append(ranking, models.Ranking{
			Position: position,
			UserID:   id,
			Score:    score,
		})
		position++
	}

	return ranking, nil
}

func (b *BasicStoreService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	user := new(models.User)
	err := b.core.DB.QueryRowContext(ctx, "SELECT id, score FROM users WHERE id = $1", id).Scan(
		&user.UserID,
		&user.Score,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
