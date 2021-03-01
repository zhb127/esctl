package up

import (
	"esctl/pkg/es"
	"esctl/pkg/log"
	tdES "esctl/test/data/pkg/es"
	tdLog "esctl/test/data/pkg/log"
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	type args struct {
		logHelper log.IHelper
		esHelper  es.IHelper
	}
	tests := []struct {
		name string
		args args
		want IService
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.args.logHelper, tt.args.esHelper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_InitMigrationHistoryRepo(t *testing.T) {
	mockLogHelper := tdLog.MockHelper()
	mockESHelper := tdES.MockHelper(mockLogHelper)

	type fields struct {
		logHelper log.IHelper
		esHelper  es.IHelper
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				logHelper: mockLogHelper,
				esHelper:  mockESHelper,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				logHelper: tt.fields.logHelper,
				esHelper:  tt.fields.esHelper,
			}
			if err := s.InitMigrationHistoryRepo(); (err != nil) != tt.wantErr {
				t.Errorf("service.InitMigrationHistoryRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_SaveMigrationHistoryEntry(t *testing.T) {
	mockLogHelper := tdLog.MockHelper()
	mockESHelper := tdES.MockHelper(mockLogHelper)

	type fields struct {
		logHelper log.IHelper
		esHelper  es.IHelper
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				logHelper: mockLogHelper,
				esHelper:  mockESHelper,
			},
			args: args{
				name: "123123123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				logHelper: tt.fields.logHelper,
				esHelper:  tt.fields.esHelper,
			}
			if err := s.SaveMigrationHistoryEntry(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("service.SaveMigrationHistoryEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_DeleteMigrationHistoryEntry(t *testing.T) {
	mockLogHelper := tdLog.MockHelper()
	mockESHelper := tdES.MockHelper(mockLogHelper)

	type fields struct {
		logHelper log.IHelper
		esHelper  es.IHelper
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				logHelper: mockLogHelper,
				esHelper:  mockESHelper,
			},
			args: args{
				name: "123123123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				logHelper: tt.fields.logHelper,
				esHelper:  tt.fields.esHelper,
			}
			if err := s.DeleteMigrationHistoryEntry(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("service.DeleteMigrationHistoryEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetUpMigrationNameLastExecuted(t *testing.T) {
	mockLogHelper := tdLog.MockHelper()
	mockESHelper := tdES.MockHelper(mockLogHelper)

	type fields struct {
		logHelper log.IHelper
		esHelper  es.IHelper
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				logHelper: mockLogHelper,
				esHelper:  mockESHelper,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				logHelper: tt.fields.logHelper,
				esHelper:  tt.fields.esHelper,
			}
			got, err := s.GetUpMigrationNameLastExecuted()
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetUpMigrationNameLastExecuted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetUpMigrationNameLastExecuted() = %v, want %v", got, tt.want)
			}
		})
	}
}
