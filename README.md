# spansql-ddlfmt

## Usage

```console
go run github.com/nametake/spansql-ddlfmt/cmd/ddlfmt@latest -w schema.sql
```


## Example

Before.

```sql
/*
CD Database
*/


-- Album Table
CREATE TABLE Album (
  ID STRING(MAX) NOT NULL, -- ID comment
  -- Title head comment
  Title STRING(MAX) NOT NULL, -- Title tail comment
  Artist STRING(MAX) NOT NULL, CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);



-- Artist Table
CREATE TABLE Artist (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);
```

After.

```sql
/*
CD Database
*/

-- Album Table
CREATE TABLE Album (
  ID STRING(MAX) NOT NULL, -- ID comment
  -- Title head comment
  Title STRING(MAX) NOT NULL, -- Title tail comment
  Artist STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);

-- Artist Table
CREATE TABLE Artist (
  ID STRING(MAX) NOT NULL,
  Name STRING(MAX) NOT NULL,
  CreatedAt TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY(ID);
```
