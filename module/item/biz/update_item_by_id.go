package biz

import (
	"context"
	"g09/common"
	"g09/module/item/model"
)

type UpdateItemStore interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error
	UpdateItems(ctx context.Context, ids []int, status string) error
	DeleteItemImage(ctx context.Context, itemId int) error
}

type updateItemBiz struct {
	store UpdateItemStore
}

func NewUpdateItemBiz(store UpdateItemStore) *updateItemBiz {
	return &updateItemBiz{store: store}
}

func (biz *updateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}

	if err := biz.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}
	return nil
}

func (biz *updateItemBiz) UpdateItemsStatus(ctx context.Context, ids []int, status string) error {
	for _, id := range ids {
		item, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
		if err != nil {
			return err
		}
		if item.Status == "Deleted" {
			return model.ErrItemIsDeleted
		}
	}
	if err := biz.store.UpdateItems(ctx, ids, status); err != nil {
		return err
	}

	return nil
}

func (biz *updateItemBiz) DeleteItemImage(ctx context.Context, id int) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}

	if err := biz.store.DeleteItemImage(ctx, id); err != nil {
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	return nil
}
