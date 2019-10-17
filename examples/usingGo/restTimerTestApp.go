package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/project-flogo/contrib/activity/log"
	"github.com/project-flogo/contrib/activity/rest"
	restTrig "github.com/project-flogo/contrib/trigger/rest"
	"github.com/project-flogo/contrib/trigger/timer"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/data/coerce"
)

func RestTimerTestApp() *api.App {
	app := api.NewApp()

	// A REST Trigger to receive HTTP message
	trg := app.NewTrigger(&restTrig.Trigger{}, &restTrig.Settings{Port: 8080})
	h, _ := trg.NewHandler(&restTrig.HandlerSettings{Method: "POST", Path: "/blah/:num"})
	h.NewAction(runActivities)

	// A Timer Trigger to send HTTP message repeatedly
	tmrTrg := app.NewTrigger(&timer.Trigger{}, nil)
	tmrHandler, _ := tmrTrg.NewHandler(&timer.HandlerSettings{StartInterval: "2s", RepeatInterval: "10s"})
	tmrHandler.NewAction(runTimerActivities)

	// A REST Activity to send data to Uri
	stng := &rest.Settings{Method: "POST", Uri: "http://localhost:8080/blah/:numID",
		Headers: map[string]string{"Accept": "application/json"}}
	restAct, _ := api.NewActivity(&rest.Activity{}, stng)

	// A log Activity for logging
	logAct, _ := api.NewActivity(&log.Activity{})

	//Store in map to avoid activity instance recreation
	activities = map[string]activity.Activity{}
	activities["log"] = logAct
	activities["rest"] = restAct

	return app
}

// Stores all the activities
var activities map[string]activity.Activity

func runActivities(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {

	// Get REST Trigger Output
	trgOut := &restTrig.Output{}
	trgOut.FromMap(inputs)

	// Coerce the required outputs to string
	content, _ := coerce.ToString(trgOut.Content)
	msg, _ := coerce.ToString(trgOut.PathParams)

	// Log the input REST pathParams
	_, err := api.EvalActivity(activities["log"], &log.Input{Message: msg})
	if err != nil {
		return nil, err
	}

	sayHello("8080")

	fmt.Println("[@8080]$ Received Data : ", content)
	response := handleInput(content)
	fmt.Println("[@8080]$ Sending Response: ", response)
	fmt.Println()

	reply := &restTrig.Reply{Code: 200, Data: response}
	return reply.ToMap(), nil
}

func sayHello(act string) {
	str := "[@" + act + "]$ Hello"
	fmt.Println(str)
}

type inputData struct {
	Name   string `json:"name"`
	Value  int    `json:"value"`
	Pin    string `json:"pin"`
	Future []int  `json:"future"`
}

func handleInput(input string) map[string]interface{} {

	var in inputData
	json.Unmarshal([]byte(input), &in)

	// if err != nil {
	// 	panic("Hello, Some problem occured during json unmarshaling")
	// }

	response := make(map[string]interface{})

	response["name"] = in.Name
	response["pin"] = in.Pin
	response["value"] = in.Value
	response["future"] = in.Future

	return response
}

func runTimerActivities(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {

	sayHello("Timer")

	go testRest()

	return nil, nil
}

func testRest() {

	input := "{\"name\": \"GPIOA\",  \"value\": 1,  \"pin\": \"PA2\", \"future\": [1, 0, 0, 1] }"
	fmt.Println("[@9096]$ Sending Data: ", input)
	fmt.Println()

	output, err := api.EvalActivity(activities["rest"],
		&rest.Input{PathParams: map[string]string{"numID": "123"},
			Content: input})

	if err != nil {
		return
	}
	fmt.Println("[@9096]$ Received Response : ", output)
	fmt.Println()
}
