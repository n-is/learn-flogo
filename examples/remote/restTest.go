package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/project-flogo/contrib/activity/log"
	restTrig "github.com/project-flogo/contrib/trigger/rest"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/stream/activity/aggregate"
)

// Stores all the activities of this app
var Activities map[string]activity.Activity

func RestTest() *api.App {
	app := api.NewApp()

	// REST Trigger to receive HTTP message
	trg := app.NewTrigger(&restTrig.Trigger{}, &restTrig.Settings{Port: 9090})
	h, _ := trg.NewHandler(&restTrig.HandlerSettings{Method: "POST", Path: "/stream"})
	h.NewAction(runActivitiesStream)

	// A log Activity for logging
	logAct, _ := api.NewActivity(&log.Activity{})

	// An Aggregate Activity to aggregate data obtained at 9090 port
	aggStng1 := &aggregate.Settings{Function: "accumulate", WindowType: "tumbling",
		WindowSize: 3, ProceedOnlyOnEmit: true}
	// addStng := map[string]string{"type": "int"}
	aggAct1, _ := api.NewActivity(&aggregate.Activity{}, aggStng1)
	aggStng2 := &aggregate.Settings{Function: "avg", WindowType: "tumbling",
		WindowSize: 5, ProceedOnlyOnEmit: false}
	aggAct2, _ := api.NewActivity(&aggregate.Activity{}, aggStng2)

	//Store in map to avoid activity instance recreation
	Activities = map[string]activity.Activity{}
	Activities["log"] = logAct
	Activities["agg1"] = aggAct1
	Activities["agg2"] = aggAct2

	return app
}

func runActivitiesStream(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {

	// Get REST Trigger Output
	trgOut := &restTrig.Output{}
	trgOut.FromMap(inputs)

	// Coerce the required outputs to string
	content, _ := coerce.ToString(trgOut.Content)

	response := handleStreamInput(content)

	reply := &restTrig.Reply{Code: 200, Data: response}
	return reply.ToMap(), nil
}

type inputStreamData struct {
	Value float64 `json:"value"`
}

func handleStreamInput(input string) map[string]interface{} {

	var in inputStreamData
	err := json.Unmarshal([]byte(input), &in)

	if err != nil {
		fmt.Println("Hello, Some problem occured during json unmarshaling")
		return nil
	}

	response := make(map[string]interface{})
	response["value"] = in.Value

	tmp := &aggregate.Input{Value: in.Value}
	output, err := api.EvalActivity(Activities["agg1"], tmp)

	if err != nil {
		return nil
	}

	if output["report"] == true {
		fmt.Println("[@9090]$ Aggregator1 Output : ", output["result"])
	}

	output, err = api.EvalActivity(Activities["agg2"], tmp)

	if err != nil {
		return nil
	}

	if output["report"] == true {
		fmt.Printf("[@9090]$ Aggregator2 Output : %0.4f\n", output["result"])
		fmt.Println()
	}

	return response
}
