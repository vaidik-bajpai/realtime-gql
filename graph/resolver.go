package graph

import "github.com/vaidik-bajpai/realtime-gql/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	subscriptionResolver *subscriptionResolver
}

func NewResolver() *Resolver {
	return &Resolver{
		subscriptionResolver: &subscriptionResolver{
			observers: make(map[string][]chan *model.Message),
		},
	}
}
