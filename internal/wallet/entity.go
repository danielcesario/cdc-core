package wallet

import (
	"github.com/danielcesario/cdc-core/internal/user"
	"gorm.io/gorm"
)

type CollaboratorType int

const (
	Admin  CollaboratorType = iota // 0
	Editor                         // 1
	Viewer                         // 2
)

type Wallet struct {
	gorm.Model
	ID            uint64 `gorm:"autoIncrement"`
	Name          string
	Code          string
	UserID        uint64
	User          user.User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Collaborators []user.User `gorm:"many2many:user_wallet;"`
	Active        bool
}

func (w *Wallet) ToResponse() *WalletResponse {
	var colaboratos []user.UserResponse
	for _, user := range w.Collaborators {
		colaboratos = append(colaboratos, *user.ToResponse())
	}

	return &WalletResponse{
		Name:          w.Name,
		Code:          w.Code,
		Collaborators: colaboratos,
	}
}

type WalletRequest struct {
	Name string `json:"name"`
}

func (wr *WalletRequest) toWallet() *Wallet {
	return &Wallet{
		Name: wr.Name,
	}
}

type WalletResponse struct {
	Name          string              `json:"name"`
	Code          string              `json:"code"`
	Collaborators []user.UserResponse `json:"collaboratos,omitempty"`
}

type WalletCollaboratorRequest struct {
	UserCode string `json:"user_code"`
	Function int    `json:"function"`
}
