package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

type mariaDBRepository struct {
	db *gorm.DB
}

func NewMariaDBRepository(db *gorm.DB) Repository {
	return &mariaDBRepository{
		db: db,
	}
}

func (r *mariaDBRepository) Store(transaction Transaction) (uint64, error) {
	result := r.db.Save(&transaction)
	if result.Error != nil {
		return 0, result.Error
	}
	return transaction.ID, nil
}

func (r *mariaDBRepository) FindByCode(code string) (*Transaction, error) {
	var transaction Transaction

	record := r.db.
		//Model(&Transaction{}).
		Preload("User").
		Preload("Wallet", func(db *gorm.DB) *gorm.DB { return db.Select("ID", "Name", "Code") }).
		Preload("PaymentMethod").
		Preload("Category").
		Preload("Entries").
		Where("code = ?", code).
		First(&transaction)

	if record.Error != nil {
		return nil, record.Error
	}

	return &transaction, nil
}

func (r *mariaDBRepository) Count(search SearchTransactionRequest) (int64, error) {
	var total int64
	r.db.Model(&Entry{}).
		Joins("left join transaction on transaction.id = entry.transaction_id").
		Joins("left join wallet on wallet_id = transaction.wallet_id and wallet.code = ?", search.Wallet).
		Where("entry.due_date between ? and ?", search.From, search.To).
		Count(&total)
	return total, nil
}

func (r *mariaDBRepository) Search(search SearchTransactionRequest) ([]TransactionSearchResult, error) {
	var results []TransactionSearchResult
	r.db.Model(&Entry{}).
		Select("entry.code, transaction.transaction_type as type, entry.amount, entry.due_date, transaction.description, user.code as user, payment_method.code as payment_method, category.code as category").
		Joins("left join transaction on transaction.id = entry.transaction_id").
		Joins("left join user on user.id = transaction.user_id").
		Joins("left join payment_method on payment_method.id = transaction.payment_method_id").
		Joins("left join category on category.id = transaction.category_id").
		Joins("left join wallet on wallet_id = transaction.wallet_id and wallet.code = ?", search.Wallet).
		Where("entry.due_date between ? and ?", search.From, search.To).
		Limit(search.Params.Limit.Size).
		Offset((search.Params.Limit.Page - 1) * search.Params.Limit.Size).
		Order(fmt.Sprintf("%s %s", search.Params.Sort.Field, search.Params.Sort.Direction)).
		Scan(&results)

	return results, nil
}
