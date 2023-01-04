package driver

import (
	_ "embed"
	"encoding/json"
	"flag"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

//go:embed test_data_model.json
var testDatamodel []byte

var flagOverwriteGolden = flag.Bool("overwrite-golden", false, "Overwrite the golden file with the current execution results")

func TestAssemble(t *testing.T) {
	var dataModel Datamodel
	err := json.Unmarshal(testDatamodel, &dataModel)
	if err != nil {
		t.Fatalf("could not decode test_data_model.json: %v", err)
	}

	tests := []struct {
		name       string
		config     Config
		goldenJson string
	}{
		{
			name: "default",
			config: Config{
				Datamodel: dataModel,
				Schema:    "public",
			},
			goldenJson: "prisma.golden.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Driver{config: tt.config}
			info, err := p.Assemble()
			if err != nil {
				t.Fatal(err)
			}

			got, err := json.MarshalIndent(info, "", "\t")
			if err != nil {
				t.Fatal(err)
			}

			if *flagOverwriteGolden {
				if err = os.WriteFile(tt.goldenJson, got, 0o664); err != nil {
					t.Fatal(err)
				}
				return
			}

			want, err := os.ReadFile(tt.goldenJson)
			if err != nil {
				t.Fatal(err)
			}

			// if diff := cmp.Diff(exp, spp); diff != "" {
			// t.Fatal(diff)
			// }
			require.JSONEq(t, string(want), string(got))
		})
	}
}