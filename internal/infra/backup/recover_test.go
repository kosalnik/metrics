package backup_test

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kosalnik/metrics/internal/infra/backup"
	"github.com/kosalnik/metrics/internal/infra/storage/mock"
	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	cases := map[string]struct {
		file       func() *os.File
		wantUpdate int
		wantErr    bool
	}{
		"No File": {
			file: func() *os.File {
				f, err := os.CreateTemp(os.TempDir(), "test")
				require.NoError(t, err)
				require.NoError(t, f.Close())
				require.NoError(t, os.Remove(f.Name()))
				return f
			},
			wantErr: true,
		},
		"No Data": {
			file: func() *os.File {
				f, err := os.CreateTemp(os.TempDir(), "test")
				require.NoError(t, err)
				_, err = f.WriteString(`{"Data":[]}`)
				require.NoError(t, err)
				require.NoError(t, f.Close())
				return f
			},
			wantUpdate: 0,
			wantErr:    false,
		},
		"With Data": {
			file: func() *os.File {
				f, err := os.CreateTemp(os.TempDir(), "test")
				require.NoError(t, err)
				_, err = f.WriteString(`{"Data":[{"ID":"a","MType":"gauge","Value":1}]}`)
				require.NoError(t, err)
				require.NoError(t, f.Close())
				return f
			},
			wantUpdate: 1,
			wantErr:    false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			s := mock.NewMockStorage(ctrl)
			s.EXPECT().UpsertAll(gomock.Any(), gomock.Any()).Times(tt.wantUpdate)

			f := tt.file()
			defer os.Remove(f.Name())

			m := backup.NewRecover(s, f.Name())
			err := m.Recover(context.Background())
			if tt.wantErr {
				require.NotEmpty(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRecover_Negative(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	_, err = f.WriteString(`{"Data":[{"ID":"a","MType":"gauge","Value":1}]}`)
	require.NoError(t, err)
	require.NoError(t, f.Close())

	cases := map[string]struct {
		path   string
		enable bool
	}{
		"No path": {path: ""},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			s := mock.NewMockStorage(ctrl)
			s.EXPECT().UpsertAll(gomock.Any(), gomock.Any()).Times(0)

			m := backup.NewRecover(s, tt.path)
			require.NoError(t, m.Recover(context.Background()))
		})
	}
}
