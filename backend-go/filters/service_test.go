package filters

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupTestService(t *testing.T) (*FilterService, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	logger, _ := zap.NewDevelopment()
	registry, _ := NewFilterRegistry([]string{"localhost:11211"})

	service := NewFilterService(db, registry, logger)

	cleanup := func() {
		db.Close()
	}

	return service, mock, cleanup
}

func TestGetFiltersByState(t *testing.T) {
	tests := []struct {
		name          string
		state         string
		mockSetup     func(sqlmock.Sqlmock)
		expectedError error
		expectedLen   int
	}{
		{
			name:  "successful query",
			state: "active",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "state", "description", "is_active"}).
					AddRow(uuid.New(), "filter1", "active", "description1", true).
					AddRow(uuid.New(), "filter2", "active", "description2", true)
				mock.ExpectQuery("SELECT (.+) FROM filters WHERE state = \\$1").
					WithArgs("active").
					WillReturnRows(rows)
			},
			expectedError: nil,
			expectedLen:   2,
		},
		{
			name:  "empty state",
			state: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// No DB calls expected
			},
			expectedError: ErrInvalidFilterState,
			expectedLen:   0,
		},
		{
			name:  "database error",
			state: "active",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM filters WHERE state = \\$1").
					WithArgs("active").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: ErrDatabaseOperation,
			expectedLen:   0,
		},
		{
			name:  "no results",
			state: "inactive",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM filters WHERE state = \\$1").
					WithArgs("inactive").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "state", "description", "is_active"}))
			},
			expectedError: nil,
			expectedLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mock, cleanup := setupTestService(t)
			defer cleanup()

			tt.mockSetup(mock)

			filters, err := service.GetFiltersByState(context.Background(), tt.state)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Len(t, filters, tt.expectedLen)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
