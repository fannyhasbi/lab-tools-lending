package service

import (
	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/repository"
	"github.com/fannyhasbi/lab-tools-lending/repository/postgres"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type BorrowService struct {
	Query      repository.BorrowQuery
	Repository repository.BorrowRepository
}

func NewBorrowService() *BorrowService {
	var borrowQuery repository.BorrowQuery
	var borrowRepository repository.BorrowRepository

	db := config.InitPostgresDB()
	borrowQuery = postgres.NewBorrowQueryPostgres(db)
	borrowRepository = postgres.NewBorrowRepositoryPostgres(db)

	return &BorrowService{
		Query:      borrowQuery,
		Repository: borrowRepository,
	}
}

func (bs BorrowService) SaveBorrow(borrow types.Borrow) (types.Borrow, error) {
	result, err := bs.Repository.Save(&borrow)
	if err != nil {
		return types.Borrow{}, err
	}

	return result, nil
}

func (bs BorrowService) UpdateBorrow(borrow types.Borrow) (types.Borrow, error) {
	result, err := bs.Repository.Update(&borrow)
	if err != nil {
		return types.Borrow{}, err
	}

	return result, nil
}

func (bs BorrowService) FindByID(id int64) ([]types.Borrow, error) {
	result := bs.Query.FindByUserID(id)
	if result.Error != nil {
		return []types.Borrow{}, result.Error
	}

	return result.Result.([]types.Borrow), nil
}
