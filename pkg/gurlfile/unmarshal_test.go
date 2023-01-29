package gurlfile

import (
	"bytes"
	"io/fs"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/gurl/mocks"
	"github.com/studio-b12/gurl/pkg/set"
)

func TestUnmarshal(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockFs := mocks.NewMockFS(mockCtl)

	contentA := `
use b/b

GET https://example1.com
	`
	mockFileB := fileMock(mockCtl, `
use ../c.gurl

GET https://example2.com
	`)
	mockFileC := fileMock(mockCtl, `
GET https://example3.com
	`)

	mockFs.EXPECT().Open("test/b/b.gurl").Return(mockFileB, nil)
	mockFs.EXPECT().Open("test/c.gurl").Return(mockFileC, nil)

	gf, err := unmarshal(mockFs, contentA, "test", set.Set[string]{})
	assert.Nil(t, err)
	assert.Equal(t, Gurlfile{
		Tests: []Request{
			testRequest("GET", "https://example3.com"),
			testRequest("GET", "https://example2.com"),
			testRequest("GET", "https://example1.com"),
		},
	}, gf)
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
