package api

import (
	"fmt"
	"io"
	"reflect"

	"github.com/ankeesler/anwork/task"
)

//go:generate go run ../cmd/genapidoc/main.go ../doc/API.md

type extraRouteData struct {
	description string
	inputType   reflect.Type
	outputType  reflect.Type
}

var erd = map[string]extraRouteData{
	"auth": extraRouteData{
		description: "create (encrypted) authentication token",
		outputType:  reflect.TypeOf(""),
	},

	"get_tasks": extraRouteData{
		description: "get all tasks",
		outputType:  reflect.SliceOf(reflect.TypeOf(task.Task{})),
	},
	"create_task": extraRouteData{
		description: "create a task",
		inputType:   reflect.TypeOf(task.Task{}),
	},
	"get_task": extraRouteData{
		description: "get a task",
		outputType:  reflect.TypeOf(task.Task{}),
	},
	"update_task": extraRouteData{
		description: "update a task",
		inputType:   reflect.TypeOf(task.Task{}),
	},
	"delete_task": extraRouteData{
		description: "delete a task",
	},

	"get_events": extraRouteData{
		description: "get all events",
		outputType:  reflect.SliceOf(reflect.TypeOf(task.Event{})),
	},
	"create_event": extraRouteData{
		description: "create an event",
		inputType:   reflect.TypeOf(task.Event{}),
	},
	"get_event": extraRouteData{
		description: "get an event",
		outputType:  reflect.TypeOf(task.Event{}),
	},
	"delete_event": extraRouteData{
		description: "delete an event",
	},
}

// MarkdownUsage will print usage documentation for the ANWORK API to an io.Writer.
// The output will be in Github markdown format.
func MarkdownUsage(output io.Writer) {
	for _, route := range routes {
		fmt.Fprintf(output, "### `%s`: `%s %s`\n", route.Name, route.Method, route.Path)

		if extra, ok := erd[route.Name]; ok {
			fmt.Fprintf(output, "* %s\n", extra.description)
			fmt.Fprintf(output, "* input: `%s`\n", typeName(extra.inputType))
			fmt.Fprintf(output, "* output: `%s`\n", typeName(extra.outputType))
		} else {
			panic(fmt.Sprintf("missing extraRouteData for '%s'", route.Name))
		}
	}
}

func typeName(teyep reflect.Type) string {
	if teyep == nil {
		return "<none>"
	} else {
		return teyep.String()
	}
}
