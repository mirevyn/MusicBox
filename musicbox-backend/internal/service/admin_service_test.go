package service

import (
	"errors"
	"testing"

	"musicbox-backend/internal/model"
)

func TestValidateAdminUserMutationRejectsProtectedAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target model.User
	}{
		{
			name:   "seed admin by id",
			target: model.User{ID: 1, Username: "root"},
		},
		{
			name:   "seed admin by username",
			target: model.User{ID: 9, Username: "admin"},
		},
		{
			name:   "seed admin by case insensitive username",
			target: model.User{ID: 9, Username: "Admin"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateAdminUserMutation(99, tt.target)
			if !errors.Is(err, ErrProtectedAdmin) {
				t.Fatalf("validateAdminUserMutation() error = %v, want ErrProtectedAdmin", err)
			}
		})
	}
}

func TestValidateAdminUserMutationRejectsSelfOperation(t *testing.T) {
	t.Parallel()

	target := model.User{ID: 7, Username: "manager"}
	err := validateAdminUserMutation(7, target)
	if !errors.Is(err, ErrSelfOperation) {
		t.Fatalf("validateAdminUserMutation() error = %v, want ErrSelfOperation", err)
	}
}

func TestValidateAdminUserMutationAllowsOtherUsers(t *testing.T) {
	t.Parallel()

	target := model.User{ID: 8, Username: "editor"}
	if err := validateAdminUserMutation(7, target); err != nil {
		t.Fatalf("validateAdminUserMutation() error = %v, want nil", err)
	}
}
