package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_MetaData_IsDeleted(t *testing.T) {
	for _, tc := range []struct {
		name string
		give MetaData
		want bool
	}{
		{
			"created",
			MetaData{
				CreatedAt: time.Now(),
			},
			false,
		},
		{
			"updated",
			MetaData{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now().Add(time.Second),
			},
			false,
		},
		{
			"deleted",
			MetaData{
				CreatedAt: time.Now(),
				DeletedAt: time.Now().Add(time.Second),
			},
			true,
		},
		{
			"deleted",
			MetaData{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now().Add(1 * time.Second),
				DeletedAt: time.Now().Add(2 * time.Second),
			},
			true,
		},
		{
			"updated after deleted",
			MetaData{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now().Add(2 * time.Second),
				DeletedAt: time.Now().Add(1 * time.Second),
			},
			false,
		},
	} {
		require.Equal(t, tc.give.IsDeleted(), tc.want)
	}
}
