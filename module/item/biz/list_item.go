package biz

import (
	"context"
	"g09/common"
	"g09/module/item/model"
)

type ListItemRepo interface {
	ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error)
}

type listItemBiz struct {
	store     ListItemRepo
	requester common.Requester
}

func NewListItemBiz(store ListItemRepo, requester common.Requester) *listItemBiz {
	return &listItemBiz{store: store, requester: requester}
}

func (biz *listItemBiz) ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging, moreKeys ...string) ([]model.TodoItem, error) {
	// ctxStore := context.WithValue(ctx, common.CurrentUser, biz.requester)
	data, err := biz.store.ListItem(ctx, filter, paging, "Owner")
	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}
	return data, nil
}
