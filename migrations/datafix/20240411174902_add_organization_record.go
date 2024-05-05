package datafix

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddOrganizationRecord, downAddOrganizationRecord)
}

func upAddOrganizationRecord(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO organizations(name) VALUES('org1');")
	return err
}

func downAddOrganizationRecord(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "ALTER TABLE organizations DISABLE TRIGGER prevent_delete_on_organizations;")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM organizations WHERE name = 'org1';")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "ALTER TABLE organizations ENABLE TRIGGER prevent_delete_on_organizations;")

	return err
}
