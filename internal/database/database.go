package database

import (
	"github.com/danielcesario/cdc-core/internal/plan"
	"github.com/danielcesario/cdc-core/internal/user"
	"github.com/danielcesario/cdc-core/internal/wallet"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func UpdateSchema(db *gorm.DB) {

	db.AutoMigrate(
		&user.Role{},
		&user.User{},
		&plan.Plan{},
		&wallet.Wallet{},
	)

	roleSuperAdmin := &user.Role{Role: "SUPER_ADMIN"}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roleSuperAdmin)

	roleUserDefault := &user.Role{Role: "USER_DEFAULT"}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roleUserDefault)
}
