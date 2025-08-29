package model

type UpdateItemsStatus struct {
	Ids    []int  `json:"ids"`
	Status string `json:"status"`
}
