package service

import (
	"fmt"
	"time"

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

func (bs BorrowService) SaveBorrow(borrow types.Borrow) (int64, error) {
	result, err := bs.Repository.Save(&borrow)
	if err != nil {
		return int64(0), err
	}

	return result, nil
}

func (bs BorrowService) UpdateBorrowStatus(id int64, status types.BorrowStatus) error {
	return bs.Repository.UpdateStatus(id, status)
}

func (bs BorrowService) UpdateBorrowConfirm(id int64, confirmedAt time.Time, firstName, lastName string) error {
	confirmedBy := firstName
	if len(lastName) > 0 {
		confirmedBy = fmt.Sprintf("%s %s", firstName, lastName)
	}

	return bs.Repository.UpdateConfirm(id, confirmedAt, confirmedBy)
}

func (bs BorrowService) FindBorrowByID(id int64) (types.Borrow, error) {
	result := bs.Query.FindByID(id)
	if result.Error != nil {
		return types.Borrow{}, result.Error
	}

	return result.Result.(types.Borrow), nil
}

func (bs BorrowService) FindByUserID(id int64) ([]types.Borrow, error) {
	result := bs.Query.FindByUserID(id)
	if result.Error != nil {
		return []types.Borrow{}, result.Error
	}

	return result.Result.([]types.Borrow), nil
}

func (bs BorrowService) GetCurrentlyBeingBorrowedByUserID(id int64) ([]types.Borrow, error) {
	status := []types.BorrowStatus{types.GetBorrowStatus("progress")}
	result := bs.Query.GetByUserIDAndMultipleStatus(id, status)
	if result.Error != nil {
		return []types.Borrow{}, result.Error
	}

	return result.Result.([]types.Borrow), nil
}

func (bs BorrowService) GetCurrentlyBeingBorrowedAndRequestedByUserID(id int64) ([]types.Borrow, error) {
	status := []types.BorrowStatus{
		types.GetBorrowStatus("request"),
		types.GetBorrowStatus("progress"),
	}
	result := bs.Query.GetByUserIDAndMultipleStatus(id, status)
	if result.Error != nil {
		return []types.Borrow{}, result.Error
	}

	return result.Result.([]types.Borrow), result.Error
}

func (bs BorrowService) GetBorrowRequests() ([]types.Borrow, error) {
	result := bs.Query.GetByStatus(types.GetBorrowStatus("request"))
	if result.Error != nil {
		return []types.Borrow{}, result.Error
	}

	return result.Result.([]types.Borrow), nil
}

func (bs BorrowService) GetBorrowReport(year, month int) ([]types.Borrow, error) {
	result := bs.Query.GetReport(year, month)
	if result.Error != nil {
		return []types.Borrow{}, result.Error
	}

	return result.Result.([]types.Borrow), nil
}
