package filters

import (
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

const (
	ivValue    = "value"
	ovFiltered = "filtered"
	ovValue    = "value"
)

type Settings struct {
	Type              string `md:"type,allowed(non-zero)"`
	ProceedOnlyOnEmit bool
}

type Input struct {
	Value interface{} `md:"value"`
}

type Output struct {
	Filtered bool        `md:"filtered"`
	Value    interface{} `md:"value"`
}

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Associates the filter name to the filter
func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{ProceedOnlyOnEmit: true}
	err := metadata.MapToStruct(ctx.Settings(), s, true)

	if err != nil {
		return nil, err
	}

	act := &Activity{}

	if s.Type == "non-zero" {
		act.filter = &NonZeroFilter{}
	} else {
		err := fmt.Errorf("Unsupported Filter : %s", s.Type)
		return nil, err
	}

	return act, nil
}

type Filter interface {
	FilterOut(val interface{}) (bool, interface{})
}

// Activity ...
// Activity is an activity that is used to filter a message to the console
type Activity struct {
	filter            Filter
	proceedOnlyOnEmit bool
}

// Metadata returns activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	filter := a.filter
	proceedOnlyOnEmit := a.proceedOnlyOnEmit

	in := ctx.GetInput(ivValue)
	filteredOut, out := filter.FilterOut(in)

	done = !(proceedOnlyOnEmit && filteredOut)
	err = ctx.SetOutput(ovFiltered, filteredOut)

	if err != nil {
		return false, err
	}

	err = ctx.SetOutput(ovValue, out)

	if err != nil {
		return false, err
	}

	return done, nil
}
