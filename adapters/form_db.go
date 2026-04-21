package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/entities"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/repositories"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
	"github.com/marlon-clemente/timyo-playground-backend/packages/database"
	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
	"gorm.io/gorm"
)

type FormDB struct {
	db *gorm.DB
}

func NewFormDB() repositories.Forms {
	return &FormDB{db: database.Get()}
}

func (r *FormDB) Count(ctx context.Context, workspaceID string) (int, error) {
	var count int64
	query := `SELECT COUNT(*) FROM forms WHERE workspace_id = ?;`
	err := r.db.WithContext(ctx).Raw(query, workspaceID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *FormDB) Create(ctx context.Context, form *entities.Form) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	insertFormSQL := `
		INSERT INTO forms (
			id,
			workspace_id,
			agent_id,
			name,
			description,
			is_public,
			current_version_id,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	err := tx.Exec(
		insertFormSQL,
		form.Id,
		form.WorkspaceID,
		form.AgentID,
		form.Name,
		form.Description,
		false,
		form.CurrentVersion,
		form.CreatedAt,
		form.UpdatedAt,
	).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	insertVersionSQL := `
		INSERT INTO form_versions (
			id,
			form_id,
			version_number,
			props,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?);
	`

	for _, version := range form.Versions {
		var props any
		if version.Props != nil {
			propsBytes, marshalErr := json.Marshal(version.Props)
			if marshalErr != nil {
				tx.Rollback()
				return marshalErr
			}
			props = propsBytes
		}

		err = tx.Exec(
			insertVersionSQL,
			version.Id,
			version.FormID,
			version.VersionNumber,
			props,
			version.CreatedAt,
			version.UpdatedAt,
		).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *FormDB) FindByID(ctx context.Context, id string) (*entities.Form, error) {
	type formRow struct {
		ID               vos.UUID        `gorm:"column:id"`
		WorkspaceID      vos.UUID        `gorm:"column:workspace_id"`
		AgentID          vos.UUID        `gorm:"column:agent_id"`
		Name             vos.Name        `gorm:"column:name"`
		Description      vos.Description `gorm:"column:description"`
		CurrentVersionID vos.UUID        `gorm:"column:current_version_id"`
		CreatedAt        time.Time       `gorm:"column:created_at"`
		UpdatedAt        time.Time       `gorm:"column:updated_at"`
	}

	query := `
		SELECT id, workspace_id, agent_id, name, description, current_version_id, created_at, updated_at
		FROM forms
		WHERE id = ?
		LIMIT 1;
	`

	var row formRow
	tx := r.db.WithContext(ctx).Raw(query, id).Scan(&row)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errs.NotFoundErr("form not found", gorm.ErrRecordNotFound).WithCode("REQUIRED_FORM")
	}

	form := &entities.Form{
		Id:             row.ID,
		WorkspaceID:    row.WorkspaceID,
		AgentID:        row.AgentID,
		Name:           row.Name,
		Description:    row.Description,
		CurrentVersion: row.CurrentVersionID,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}

	return form, nil
}

func (r *FormDB) ListByWorkspaceID(ctx context.Context, workspaceID vos.UUID) ([]*entities.Form, error) {
	type formRow struct {
		ID               vos.UUID        `gorm:"column:id"`
		WorkspaceID      vos.UUID        `gorm:"column:workspace_id"`
		AgentID          vos.UUID        `gorm:"column:agent_id"`
		Name             vos.Name        `gorm:"column:name"`
		Description      vos.Description `gorm:"column:description"`
		CurrentVersionID vos.UUID        `gorm:"column:current_version_id"`
		CreatedAt        time.Time       `gorm:"column:created_at"`
		UpdatedAt        time.Time       `gorm:"column:updated_at"`
	}

	query := `
		SELECT id, workspace_id, agent_id, name, description, current_version_id, created_at, updated_at
		FROM forms
		WHERE workspace_id = ?;
	`

	var rows []formRow
	err := r.db.WithContext(ctx).Raw(query, workspaceID.String()).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	var forms []*entities.Form
	for _, row := range rows {
		
		fmt.Printf("row: %v\n", row)

		form := &entities.Form{
			Id:             row.ID,
			WorkspaceID:    row.WorkspaceID,
			AgentID:        row.AgentID,
			Name:           row.Name,
			Description:    row.Description,
			CurrentVersion: row.CurrentVersionID,
			CreatedAt:      row.CreatedAt,
			UpdatedAt:      row.UpdatedAt,
		}
		forms = append(forms, form)
	}

	return forms, nil
}

func (r *FormDB) GetByIDWithVersion(ctx context.Context, id vos.UUID) (*entities.Form, error) {
	form, err := r.FindByID(ctx, id.String())
	if err != nil {
		return nil, err
	}

	type versionRow struct {
		ID            vos.UUID  `gorm:"column:id"`
		FormID        vos.UUID  `gorm:"column:form_id"`
		VersionNumber int       `gorm:"column:version_number"`
		Props         []byte    `gorm:"column:props"`
		CreatedAt     time.Time `gorm:"column:created_at"`
		UpdatedAt     time.Time `gorm:"column:updated_at"`
	}

	query := `
		SELECT id, form_id, version_number, props, created_at, updated_at
		FROM form_versions
		WHERE form_id = ?;
	`

	var rows []versionRow
	err = r.db.WithContext(ctx).Raw(query, id.String()).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		var props any
		if len(row.Props) > 0 {
			err = json.Unmarshal(row.Props, &props)
			if err != nil {
				return nil, err
			}
		}

		version := &entities.FormVersion{
			Id:            row.ID,
			FormID:        row.FormID,
			VersionNumber: row.VersionNumber,
			Props:         props,
			CreatedAt:     row.CreatedAt,
			UpdatedAt:     row.UpdatedAt,
		}
		form.Versions = append(form.Versions, *version)
	}

	return form, nil
}

func (r *FormDB) UpdateFormVersion(ctx context.Context, version *entities.FormVersion) error {
	var props any
	if version.Props != nil {
		propsBytes, marshalErr := json.Marshal(version.Props)
		if marshalErr != nil {
			return marshalErr
		}
		props = propsBytes
	}

	updateSQL := `
		UPDATE form_versions
		SET props = ?, updated_at = ?
		WHERE id = ?;
	`

	return r.db.WithContext(ctx).Exec(updateSQL, props, version.UpdatedAt, version.Id).Error
}

func (r *FormDB) Delete(ctx context.Context, workspaceID, formID vos.UUID) error {
	deleteFormSQL := `DELETE FROM forms WHERE id = ? AND workspace_id = ?;`

	tx := r.db.WithContext(ctx).Exec(deleteFormSQL, formID.String(), workspaceID.String())
	if tx.Error != nil {
		return fmt.Errorf("failed to delete form: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("form not found: %w", gorm.ErrRecordNotFound)
	}

	return nil
}