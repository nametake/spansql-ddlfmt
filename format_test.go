package ddlfmt

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
			ddl: `CREATE TABLE Album (
ID STRING(MAX) NOT NULL,
Title STRING(MAX) NOT NULL,
Artist STRING(MAX) NOT NULL, CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
			want: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL,
  Title STRING(MAX) NOT NULL,
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
		},
		{
			name: "Two DDL",
			ddl: `CREATE TABLE Album (
ID STRING(MAX) NOT NULL,
Title STRING(MAX) NOT NULL,
Artist STRING(MAX) NOT NULL,
CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);
CREATE TABLE Artist (
    ID STRING(MAX) NOT NULL, Name STRING(MAX) NOT NULL,
    CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
			want: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL,
  Title STRING(MAX) NOT NULL,
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

CREATE TABLE Artist (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
		},
		{
			name: "Stream DDL",
			ddl: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL,
  Title STRING(MAX) NOT NULL,
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

CREATE CHANGE STREAM SingerAlbumStream
FOR Singers, Albums;`,
			want: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL,
  Title STRING(MAX) NOT NULL,
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

CREATE CHANGE STREAM SingerAlbumStream FOR Singers, Albums;`,
		},
		{
			name: "Comment block",
			ddl: `CREATE TABLE Album (ID STRING(MAX) NOT NULL, Title STRING(MAX) NOT NULL, Artist STRING(MAX) NOT NULL, CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true)) PRIMARY KEY(ID);
-- COMMENT 1
-- COMMENT 2
CREATE CHANGE STREAM SingerAlbumStream FOR Singers, Albums;
			`,
			want: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL,
  Title STRING(MAX) NOT NULL,
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

-- COMMENT 1
-- COMMENT 2
CREATE CHANGE STREAM SingerAlbumStream FOR Singers, Albums;`,
		},
		{
			name: "Inline comment",
			ddl: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL, -- ID comment
  -- Title head comment
  Title STRING(MAX) NOT NULL, -- Title tail comment
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

-- COMMENT 1
-- COMMENT 2

-- Artist Table
CREATE TABLE Artist (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
			want: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL, -- ID comment
  -- Title head comment
  Title STRING(MAX) NOT NULL, -- Title tail comment
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

-- COMMENT 1
-- COMMENT 2

-- Artist Table
CREATE TABLE Artist (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
		},
		{
			name: "Comment pattern",
			ddl: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL, # ID comment
  # Title head comment
  Title STRING(MAX) NOT NULL, ## Title tail comment
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

####################
## Styled Comment ##
####################

--------------------
-- Styled Comment --
--------------------

/*
star
multi
line
comment
*/

/*
star
multi
line
leading
comment
*/
CREATE TABLE Artist (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
			want: `CREATE TABLE Album (
  ID STRING(MAX) NOT NULL, # ID comment
  # Title head comment
  Title STRING(MAX) NOT NULL, ## Title tail comment
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

####################
## Styled Comment ##
####################

--------------------
-- Styled Comment --
--------------------

/*
star
multi
line
comment
*/

/*
star
multi
line
leading
comment
*/
CREATE TABLE Artist (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatDDL("test.sql", tt.ddl)
			if err != nil {
				t.Fatalf("FormatDDL() error = %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FormatDDL() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
