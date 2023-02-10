package api

import "github.com/gin-gonic/gin"

// Action Abstract strategy for action
type Action interface {
	Do(c *gin.Context) (data interface{})
}

// Actions action factory
var Actions = map[string]Action{
	"EmptyAction": new(EmptyAction),
}

// EmptyAction Just for example
type EmptyAction struct{}

func (e EmptyAction) Do(c *gin.Context) (data interface{}) {
	// TODO implement me
	panic("implement me")
}
