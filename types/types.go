package types

import (
	"github.com/smolgu/katalog/models/user"
)

// Store represent db methods
type Store interface {
	// training directions
	//GetTrainingDirection([]byte) (*user.TrainingDirection, error)
	GetTrainingDirectionByShortName(string) (*user.TrainingDirection, error)
	SetTrainingDirection(*user.TrainingDirection) error

	UserSave(*user.User) error
	UsersGet([]byte) (*user.User, error)
	UsersGetByVkID(int64) (*user.User, error)
}
