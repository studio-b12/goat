package goatfile

import (
	"bytes"
	"errors"
	"io/fs"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/goat/mocks"
	"github.com/studio-b12/goat/pkg/set"
)

func TestUnmarshal(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockFs := mocks.NewMockFS(mockCtl)

	contentA := `
use b/b

GET https://example1.com
	`
	mockFileB := fileMock(mockCtl, `
use ../c.goat

GET https://example2.com
	`)
	mockFileC := fileMock(mockCtl, `
GET https://example3.com
	`)

	mockFs.EXPECT().Open("test/b/b.goat").Return(mockFileB, nil)
	mockFs.EXPECT().Open("test/c.goat").Return(mockFileC, nil)

	gf, err := unmarshal(mockFs, contentA, "test/test.goat", set.Set[string]{})
	assert.Nil(t, err)
	assert.Equal(t, Goatfile{
		Tests: []Action{
			testRequestWithPath("GET", "https://example3.com", "test/c.goat", 2),
			testRequestWithPath("GET", "https://example2.com", "test/b/b.goat", 4),
			testRequestWithPath("GET", "https://example1.com", "test/test.goat", 4),
		},
		Path: "test/test.goat",
	}, gf)
}

func TestUnmarshal_Errors(t *testing.T) {
	t.Run("notfound", func(t *testing.T) {
		mockCtl := gomock.NewController(t)
		mockFs := mocks.NewMockFS(mockCtl)

		contentA := `
use b/b

GET https://example1.com
	`
		mockFileB := fileMock(mockCtl, `
use ../c.goat

GET https://example2.com
	`)

		errNotFound := errors.New("err not found")

		mockFs.EXPECT().Open("test/b/b.goat").Return(mockFileB, nil)
		mockFs.EXPECT().Open("test/c.goat").Return(nil, errNotFound)

		_, err := unmarshal(mockFs, contentA, "test/test.goat", set.Set[string]{})
		assert.ErrorIs(t, err, errNotFound, err)
	})

	t.Run("notfound", func(t *testing.T) {
		mockCtl := gomock.NewController(t)
		mockFs := mocks.NewMockFS(mockCtl)

		errReadError := errors.New("read error")

		contentA := `
use b/b

GET https://example1.com
	`
		mockFileB := fileMock(mockCtl, `
use ../c.goat

GET https://example2.com
	`)
		mockFileC := fileMockErr(mockCtl, errReadError)

		mockFs.EXPECT().Open("test/b/b.goat").Return(mockFileB, nil)
		mockFs.EXPECT().Open("test/c.goat").Return(mockFileC, nil)

		_, err := unmarshal(mockFs, contentA, "test/test.goat", set.Set[string]{})
		assert.ErrorIs(t, err, errReadError, err)
	})
}

func TestUnmarshal_MultipleImports(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockFs := mocks.NewMockFS(mockCtl)

	contentA := `
use b/b

GET https://example1.com
	`
	mockFileA := fileMock(mockCtl, contentA)
	mockFileB := fileMock(mockCtl, `
use ../c.goat

GET https://example2.com
	`)
	mockFileC := fileMock(mockCtl, `
use a

GET https://example3.com
	`)

	mockFs.EXPECT().Open("test/a.goat").Return(mockFileA, nil)
	mockFs.EXPECT().Open("test/b/b.goat").Return(mockFileB, nil)
	mockFs.EXPECT().Open("test/c.goat").Return(mockFileC, nil)

	_, err := unmarshal(mockFs, contentA, "test/a.goat", set.Set[string]{})
	assert.ErrorIs(t, err, ErrMultiImport)
}

// --- Helpers ---

func fileMock(mockCtl *gomock.Controller, raw string) fs.File {
	mockFile := mocks.NewMockFile(mockCtl)
	buf := bytes.NewBufferString(raw)

	mockFileInfo := mocks.NewMockFileInfo(mockCtl)
	mockFileInfo.EXPECT().Size().Return(int64(len(raw))).AnyTimes()

	mockFile.EXPECT().Stat().Return(mockFileInfo, nil).AnyTimes()

	mockFile.EXPECT().Read(gomock.Any()).DoAndReturn(func(v []byte) (int, error) {
		return buf.Read(v)
	}).AnyTimes()

	mockFile.EXPECT().Close().Return(nil).AnyTimes()

	return mockFile
}

func fileMockErr(mockCtl *gomock.Controller, err error) fs.File {
	mockFile := mocks.NewMockFile(mockCtl)

	mockFileInfo := mocks.NewMockFileInfo(mockCtl)
	mockFileInfo.EXPECT().Size().Return(int64(0)).AnyTimes()

	mockFile.EXPECT().Stat().Return(mockFileInfo, nil).AnyTimes()

	mockFile.EXPECT().Read(gomock.Any()).DoAndReturn(func(v []byte) (int, error) {
		return 0, err
	}).AnyTimes()

	mockFile.EXPECT().Close().Return(nil).AnyTimes()

	return mockFile
}
