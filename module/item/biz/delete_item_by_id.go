package biz

import (
	"context"
	"g09/common"
	"g09/module/item/model"
)

type DeleteItemStore interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error
	DeleteItems(ctx context.Context, ids []int) error
}

type deleteItemBiz struct {
	store DeleteItemStore
}

func NewDeleteItemBiz(store DeleteItemStore) *deleteItemBiz {
	return &deleteItemBiz{store: store}
}

func (biz *deleteItemBiz) DeleteItemById(ctx context.Context, id int) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}
	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}

	if err := biz.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return common.ErrCannotDeleteEntity(model.EntityName, err)
	}
	return nil
}

func (biz *deleteItemBiz) DeleteItemsByIds(ctx context.Context, ids []int) error {
	for _, id := range ids {
		item, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
		if err != nil {
			return common.ErrCannotGetEntity(model.EntityName, err)
		}
		if item.Status == "Deleted" {
			return model.ErrItemIsDeleted
		}
	}
	return biz.store.DeleteItems(ctx, ids)
}
