package user

import (
	"database/sql"
	"go-clean/src/business/entity"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_user_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "INSERT INTO"
	query := regexp.QuoteMeta(querySql)

	mockUser := entity.User{
		Username: "mail",
		Password: "mail",
		Name:     "mail",
		IsAdmin:  true,
	}

	type args struct {
		user entity.User
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.User
		wantErr     bool
	}{
		{
			name: "failed to create user",
			args: args{
				user: mockUser,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    mockUser,
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				user: mockUser,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()
				sqlMock.ExpectationsWereMet()
				return sqlServer, err
			},
			want:    mockUser,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			_, err = u.Create(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_user_GetByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1"
	query := regexp.QuoteMeta(querySql)

	type args struct {
		username string
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.User
		wantErr     bool
	}{
		{
			name: "failed to query row",
			args: args{
				username: "mail",
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				username: "mail",
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlMock.NewRows([]string{"username"})
				row.AddRow("mail")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			want: entity.User{
				Username: "mail",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			got, err := u.GetByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Getlist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_user_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	querySql := "SELECT * FROM `users` WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1"
	query := regexp.QuoteMeta(querySql)

	type args struct {
		id uint
	}
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.User
		wantErr     bool
	}{
		{
			name: "failed to query row",
			args: args{
				id: 1,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				id: 1,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlMock.NewRows([]string{"username"})
				row.AddRow("mail")
				sqlMock.ExpectQuery(query).WillReturnRows(row)
				return sqlServer, err
			},
			want: entity.User{
				Username: "mail",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlServer,
				SkipInitializeWithVersion: true,
			}))
			if err != nil {
				t.Error(err)
			}

			u := Init(sqlClient)
			got, err := u.GetById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.GetById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
