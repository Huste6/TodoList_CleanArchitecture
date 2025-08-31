package biz

import (
	"context"
	"g09/common"
	"g09/module/user/model"
)

type RegisterStorage interface {
	FindUser(ctx context.Context, cond map[string]interface{}, moreInfo ...string) (*model.User, error)
	CreateUser(ctx context.Context, data *model.UserCreate) error
}

type Hasher interface {
	Hash(data string) string
}

type registerBusiness struct {
	registerStorage RegisterStorage
	hasher          Hasher
}

func NewRegisterBusiness(registerStorage RegisterStorage, hasher Hasher) *registerBusiness {
	return &registerBusiness{
		registerStorage: registerStorage,
		hasher:          hasher,
	}
}

func (business *registerBusiness) Register(ctx context.Context, data *model.UserCreate) error {
	user, err := business.registerStorage.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		return model.ErrEmailExisted
	}

	if err != common.RecordNotFound {
		return err
	}

	salt := common.GenSalt(50)

	hashedPassword := business.hasher.Hash(data.Password + salt)

	data.Password = hashedPassword
	data.Salt = salt
	data.Role = model.RoleUser

	if err := business.registerStorage.CreateUser(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
