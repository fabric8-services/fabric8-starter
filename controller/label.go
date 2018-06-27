package controller

import (
	"github.com/fabric8-services/fabric8-starter/app"
	"github.com/goadesign/goa"
)

// LabelController implements the label resource.
type LabelController struct {
	*goa.Controller
}

// NewLabelController creates a label controller.
func NewLabelController(service *goa.Service) *LabelController {
	return &LabelController{Controller: service.NewController("LabelController")}
}

// Create runs the create action.
func (c *LabelController) Create(ctx *app.CreateLabelContext) error {
	// LabelController_Create: start_implement

	// Put your logic here

	// LabelController_Create: end_implement
	return nil
}

// List runs the list action.
func (c *LabelController) List(ctx *app.ListLabelContext) error {
	// LabelController_List: start_implement

	// Put your logic here

	// LabelController_List: end_implement
	res := &app.LabelList{}
	return ctx.OK(res)
}

// Show runs the show action.
func (c *LabelController) Show(ctx *app.ShowLabelContext) error {
	// LabelController_Show: start_implement

	// Put your logic here

	// LabelController_Show: end_implement
	res := &app.LabelSingle{}
	return ctx.OK(res)
}

// Update runs the update action.
func (c *LabelController) Update(ctx *app.UpdateLabelContext) error {
	// LabelController_Update: start_implement

	// Put your logic here

	// LabelController_Update: end_implement
	res := &app.LabelSingle{}
	return ctx.OK(res)
}
