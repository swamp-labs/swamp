package config

import (
	"github.com/swamp-labs/swamp/engine/assertion"
	"github.com/swamp-labs/swamp/engine/httpreq"
	"github.com/swamp-labs/swamp/engine/simulation"
	"github.com/swamp-labs/swamp/engine/volume"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
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
			args: args{filepath.Join(dir, "../example/example.yaml")},
			want: simulation.MakeSimulation(
				[]volume.Volume{},
				map[string][]httpreq.Request{
					"g1": {{
						Name:     "req1",
						Method:   "POST",
						Protocol: "https",
						Headers: []map[string]string{
							{"Content-Type": "application/json"},
							{"Accept": "*/*"},
						},
						URL:             "https://reqres.in/api/users",
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
								{"Content-Type": {"application/json; charset=utf-8"}},
							},
						),
					},
						{
							Name:   "req2",
							Method: "GET",
							Headers: []map[string]string{
								{"Content-Type": "application/json"},
								{"Accept": "*/*"},
							},
							URL:             "https://reqres.in/api/users/${id}",
							QueryParameters: nil,
							Assertions: assertion.MakeRequestAssertion(
								[]assertion.BodyAssertion{
									assertion.NewJsonPathAssertion("$.id", "id"),
								},
								nil,
								nil,
							),
						},
					},
					"g2": {{
						Name:            "req3",
						Method:          "GET",
						Headers:         nil,
						URL:             "https://reqres.in/api/users",
						QueryParameters: nil,
						Assertions: assertion.MakeRequestAssertion(
							[]assertion.BodyAssertion{},
							[]int{200},
							nil,
						),
					}},
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
			if !reflect.DeepEqual(got.GetGroups(), tt.want.GetGroups()) {
				t.Errorf("Parse() got = %v, want %v", got.GetGroups(), tt.want.GetGroups())
			}
		})
	}
}
