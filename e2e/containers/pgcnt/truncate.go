package pgcnt

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const truncateFunc = `
CREATE OR REPLACE FUNCTION f_truncate_tables(_username text)
  RETURNS void
  LANGUAGE plpgsql AS
$func$
BEGIN
  EXECUTE
  (SELECT 'TRUNCATE TABLE '
       || string_agg(format('%I.%I', schemaname, tablename), ', ')
   FROM   pg_tables
   WHERE  tableowner = _username
   AND    schemaname = 'public'
   AND    tablename != 'schema_migrations'
   AND    tablename != 'plant_category'
   );
END
$func$;
`

func truncateCommand(username string) string {
	return fmt.Sprintf("SELECT f_truncate_tables('%s')", username)
}

func TruncateTables(ctx context.Context, db *PostgresConfig) error {
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", db.User, db.Password, *db.OuterHost, *db.OuterPort, db.Database))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, truncateFunc)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, truncateCommand(db.User))
	return err
}
