package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFormatDDL(t *testing.T) {
	tests := []struct {
		name string
		ddl  string
		want string
	}{
		{
			name: "Single DDL",
			ddl:  `CREATE TABLE User (ID STRING(MAX) NOT NULL, Name STRING(MAX), Age INT64, CreatedAt TIMESTAMP NOT NULL OPTIONS ( allow_commit_timestamp = true )) PRIMARY KEY (ID);`,
			want: `CREATE TABLE User (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX),
  Age INT64,
  CreatedAt TIMESTAMP NOT NULL OPTIONS ( allow_commit_timestamp = true ),
) PRIMARY KEY (ID);
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatDDL(tt.ddl)
			if err != nil {
				t.Fatalf("FormatDDL() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FormatDDL() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
