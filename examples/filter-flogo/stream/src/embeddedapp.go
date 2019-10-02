// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

// embedded flogo app descriptor file
const flogoJSON string = `{
    "name": "stream",
    "type": "flogo:app",
    "version": "0.0.1",
    "description": "",
    "appModel": "1.1.0",
    "imports": [
        "github.com/n-is/learn-flogo/stream/activity/filters",
        "github.com/project-flogo/contrib/trigger/rest",
        "github.com/project-flogo/stream",
        "github.com/project-flogo/contrib/activity/log"
    ],
    "triggers": [
        {
            "id": "receive_http_message",
            "ref": "#rest",
            "settings": {
                "port": "9233"
            },
            "handlers": [
                {
                    "settings": {
                        "method": "GET",
                        "path": "/filter"
                    },
                    "actions": [
                        {
                            "id": "simple_agg",
                            "input": {
                                "input": "=$.content.value"
                            }
                        }
                    ]
                }
            ]
        }
    ],
    "resources": [
        {
            "id": "pipeline:simple_filter",
            "data": {
                "metadata": {
                    "input": [
                        {
                            "name": "input",
                            "type": "integer"
                        }
                    ]
                },
                "stages": [
                    {
                        "ref": "#filters",
                        "settings": {
                            "type": "non-zero",
                            "proceedOnlyOnEmit": true
                        },
                        "input": {
                            "value": "=$.input"
                        }
                    },
                    {
                        "ref": "#log",
                        "input": {
                            "message": "=$.ovalue"
                        }
                    }
                ]
            }
        }
    ],
    "actions": [
        {
            "ref": "#stream",
            "settings": {
                "pipelineURI": "res://pipeline:simple_filter"
            },
            "id": "simple_agg"
        }
    ]
}`

func init () {
	cfgJson = flogoJSON
}
