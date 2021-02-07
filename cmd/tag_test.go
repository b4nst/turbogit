package cmd

import (
	"errors"
	"testing"

	"github.com/b4nst/turbogit/internal/format"
	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
)

func TestParseDescription(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		input  string
		err    error
		v      *semver.Version
		offset int
	}{
		{"parse desc 1", "v1.0.0", nil, &[]semver.Version{semver.MustParse("1.0.0")}[0], 1},
		{"parse desc 2", "v1.0.0-2-ab23e5f1", nil, &[]semver.Version{semver.MustParse("1.0.0")}[0], 3},
		{"parse desc 3", "latest", nil, nil, 1},
		{"parse desc 4", "latest-3-5570541a", nil, nil, 4},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v, offset, err := parseDescription(tt.input)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.v, v)
			assert.Equal(t, tt.offset, offset)
		})
	}
}

func TestBumpVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		curr     *semver.Version
		bump     format.Bump
		err      error
		expected *semver.Version
	}{
		{"bump version 1", &semver.Version{Major: 0, Minor: 0, Patch: 0}, format.BUMP_PATCH, nil, &semver.Version{Major: 0, Minor: 0, Patch: 1}},
		{"bump version 2", &semver.Version{Major: 3, Minor: 4, Patch: 7}, format.BUMP_PATCH, nil, &semver.Version{Major: 3, Minor: 4, Patch: 8}},
		{"bump version 3", &semver.Version{Major: 0, Minor: 0, Patch: 0}, format.BUMP_MINOR, nil, &semver.Version{Major: 0, Minor: 1, Patch: 0}},
		{"bump version 4", &semver.Version{Major: 3, Minor: 4, Patch: 7}, format.BUMP_MINOR, nil, &semver.Version{Major: 3, Minor: 5, Patch: 0}},
		{"bump version 5", &semver.Version{Major: 0, Minor: 0, Patch: 0}, format.BUMP_MAJOR, nil, &semver.Version{Major: 0, Minor: 1, Patch: 0}},
		{"bump version 6", &semver.Version{Major: 3, Minor: 4, Patch: 7}, format.BUMP_MAJOR, nil, &semver.Version{Major: 4, Minor: 0, Patch: 0}},
		{"bump version 7", nil, format.BUMP_PATCH, errors.New("Received nil pointer"), nil},
		{"bump version 8", nil, format.BUMP_MINOR, errors.New("Received nil pointer"), nil},
		{"bump version 9", nil, format.BUMP_MAJOR, errors.New("Received nil pointer"), nil},
		{"bump version 10", &semver.Version{Major: 0, Minor: 0, Patch: 0}, format.BUMP_NONE, nil, &semver.Version{Major: 0, Minor: 0, Patch: 0}},
		{"bump version 11", &semver.Version{Major: 3, Minor: 4, Patch: 7}, format.BUMP_NONE, nil, &semver.Version{Major: 3, Minor: 4, Patch: 7}},
		{"bump version 12", nil, format.BUMP_MAJOR, errors.New("Received nil pointer"), nil},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := bumpVersion(tt.curr, tt.bump)
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, tt.curr)
		})
	}
}
