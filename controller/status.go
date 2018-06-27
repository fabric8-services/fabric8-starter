package controller

import (
	"github.com/fabric8-services/fabric8-starter/app"
	"github.com/goadesign/goa"
)

// StatusController implements the status resource.
type StatusController struct {
	*goa.Controller
}

// NewStatusController creates a status controller.
func NewStatusController(service *goa.Service) *StatusController {
	return &StatusController{Controller: service.NewController("StatusController")}
}

// Show runs the show action.
func (c *StatusController) Show(ctx *app.ShowStatusContext) error {
	// StatusController_Show: start_implement

	// Put your logic here

	// StatusController_Show: end_implement
	res := &app.Status{}
	return ctx.OK(res)
}
