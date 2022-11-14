package config

import (
	"github.com/swamp-labs/swamp/engine/assertion"
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/simulation"
	"github.com/swamp-labs/swamp/engine/task"
	"github.com/swamp-labs/swamp/engine/templateString"
	"github.com/swamp-labs/swamp/engine/volume"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

const (
	ContentTypeHeader = "Content-Type"
)

func TestParse(t *testing.T) {

	_, fileName, _, _ := runtime.Caller(0)
	dir := filepath.Dir(fileName)

	type args struct {
		configFile string
	}
	tests := []struct {
		name    string
		args    args
		want    simulation.Simulation
		wantErr bool
	}{
		{
			name: "Example file",
			args: args{filepath.Join(dir, "../example/basic.yaml")},
			want: simulation.MakeSimulation(
				map[string]task.Task{
					"t1": task.MakeTask([]httpreq.Request{
						{
							Name:     "req1",
							Method:   "POST",
							Protocol: "https",
							Headers: []map[string]string{
								{ContentTypeHeader: "application/json"},
								{"Accept": "*/*"},
							},
							URL:             templateString.TemplateString{Format: "https://reqres.in/api/users"},
							Body:            "{ \"name\": \"batman\", \"job\": \"superhero\"}",
							QueryParameters: nil,
							Assertions: assertion.MakeRequestAssertion(
								[]assertion.BodyAssertion{
									assertion.NewJsonPathAssertion("$.id", "id"),
									assertion.NewRegexAssertion("(\\d{4,})", "date"),
								},
								[]int{201},
								[]map[string][]string{
									{"Access-Control-Allow-Origin": {"*"}},
									{ContentTypeHeader: {"application/json; charset=utf-8"}},
								},
							),
						},
					},
						volume.Volume{{"wait": 1}, {"rps": 10, "during": 60}}),
				},
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.GetTasks(), tt.want.GetTasks()) {
				t.Errorf("Parse() got = %v, want %v", got.GetTasks(), tt.want.GetTasks())
			}
		})
	}
}
