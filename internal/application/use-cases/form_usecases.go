package usecases

import (
	"context"
	"errors"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/entities"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/repositories"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
)

var LimitFormsPerWorkspace = 5

type Form struct {
	repo repositories.Forms
}

func NewFormUseCase(repo repositories.Forms) *Form {
	return &Form{repo: repo}
}

func (f *Form) Create(ctx context.Context, input CreateFormInput) (*FormOutput, error) {
	
	workspaceID, err := vos.ParseUUID(input.WorkspaceID)
	if err != nil {
		return nil, errs.DomainErr("invalid workspace ID", err)
	}

	agentID, err := vos.ParseUUID(input.AgentID)
	if err != nil {
		return nil, errs.DomainErr("invalid agent ID", err)
	}

	countForms, err := f.repo.Count(ctx, workspaceID.String())
	if err != nil {
		return nil, err
	}

	name, err := vos.NewName(input.Name)
	if err != nil {
		return nil, errs.DomainErr("invalid form name", err)
	}
	
	description, err := vos.NewDescription(input.Description)
	if err != nil {
		return nil, errs.DomainErr("invalid form description", err)
	}

	if countForms >= LimitFormsPerWorkspace {
		return nil, errs.DomainErr("workspace has reached the maximum number of forms allowed", errors.New("workspace has reached the maximum number of forms allowed"))
	}

	form := entities.NewForm(workspaceID, agentID, name, description)

	err = f.repo.Create(ctx, form)
	if err != nil {
		return nil, err
	}

	return &FormOutput{
		FormID: form.Id.String(),
		FormVersionID: string(form.CurrentVersion),
		FormName: form.Name.String(),
		FormDescription: form.Description.String(),
		CreatedAt: form.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: form.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (f *Form) ListAlls(ctx context.Context, input ListFormsInput) (*ListFormsOutput, error) {
	workspaceUUID, err := vos.ParseUUID(input.WorkspaceID)
	if err != nil {
		return nil, errs.DomainErr("invalid workspace ID", err)
	}

	forms, err := f.repo.ListByWorkspaceID(ctx, workspaceUUID)
	if err != nil {
		return nil, err
	}

	var result = make([]FormOutput, 0, len(forms))
	for _, form := range forms {
		result = append(result, FormOutput{
			FormID:          form.Id.String(),
			FormName:        form.Name.String(),
			FormDescription: form.Description.String(),
			FormVersionID: string(form.CurrentVersion),
			CreatedAt: form.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: form.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &ListFormsOutput{
		Forms: result,
		Total: len(result),
	}, nil
}

func (f *Form) GetByID(ctx context.Context, formID string) (*FormOutput, error) {
	formUUID, err := vos.ParseUUID(formID)
	if err != nil {
		return nil, errs.DomainErr("invalid form ID", err)
	}

	form, err := f.repo.GetByIDWithVersion(ctx, formUUID)
	if err != nil {
		return nil, err
	}

	return &FormOutput{
		FormID:          form.Id.String(),
		FormName:        form.Name.String(),
		FormDescription: form.Description.String(),
		FormVersionID: string(form.CurrentVersion),
		CreatedAt: form.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: form.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Props: form.Versions[0].Props,
	}, nil
}

func (f *Form) SaveFormVersion(ctx context.Context, input SaveFormVersionInput) (*FormOutput, error) {
	formUUID, err := vos.ParseUUID(input.FormID)
	if err != nil {
		return nil, errs.DomainErr("invalid form ID", err)
	}

	form, err := f.repo.GetByIDWithVersion(ctx, formUUID)
	if err != nil {
		return nil, err
	}

	version := form.Versions[0]
	version.Props = input.Props

	err = f.repo.UpdateFormVersion(ctx, &version)
	if err != nil {
		return nil, err
	}

	return &FormOutput{
		FormID:          version.FormID.String(),
		FormVersionID: version.Id.String(),
		CreatedAt: version.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: version.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (f *Form) Delete(ctx context.Context, cmd DeleteFormInput) error {
	formUUID, err := vos.ParseUUID(cmd.FormID)
	if err != nil {
		return errs.DomainErr("invalid form ID", err)
	}

	workspaceUUID, err := vos.ParseUUID(cmd.WorkspaceID)
	if err != nil {
		return errs.DomainErr("invalid workspace ID", err)
	}

	err = f.repo.Delete(ctx, workspaceUUID, formUUID)
	if err != nil {
		return err
	}

	return nil
}