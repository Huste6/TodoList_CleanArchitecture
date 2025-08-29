package biz

import (
	"context"
	"g09/common"
	"g09/module/item/model"
)

// Handler -> Biz (logic biz) [-> Repository (Tổng hợp dữ liệu và transform sang từng cấu trúc mong muốn)] -> Storage (data layer)

type CreateItemStorage interface {
	CreateItem(ctx context.Context, data *model.TodoItemCreation) error
}

type createItemBiz struct {
	store CreateItemStorage
}

func NewCreateItemBiz(store CreateItemStorage) *createItemBiz {
	return &createItemBiz{store: store}
}

func (biz *createItemBiz) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := data.Validate(); err != nil {
		return common.ErrTitleEmpty(err)
	}

	if err := biz.store.CreateItem(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
