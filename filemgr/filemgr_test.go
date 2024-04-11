package filemgr

import (
	"io/fs"
	"reflect"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

var testfs = fstest.MapFS{
	"dir1": &fstest.MapFile{
		Mode: fs.ModeDir,
	},
	"dir2": &fstest.MapFile{
		Mode: fs.ModeDir,
	},
	"file1.txt": &fstest.MapFile{
		Data: []byte("file1"),
	},
	"file2.txt": &fstest.MapFile{
		Data: []byte("file2"),
	},
	"file3.txt": &fstest.MapFile{
		Data: []byte("file3"),
	},
	"binary1.bin": &fstest.MapFile{
		Data: []byte{0x01, 0x02, 0x03},
	},
	"binary2.bin": &fstest.MapFile{
		Data: []byte{0x04, 0x05, 0x06},
	},
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Test_collectFiles(t *testing.T) {
	type args struct {
		fsys  fs.FS
		globs []string
	}
	tests := []struct {
		name      string
		args      args
		wantFiles []fs.FileInfo
		wantErr   bool
	}{
		{
			name: "collect all files",
			args: args{
				fsys:  testfs,
				globs: []string{"*"},
			},
			wantFiles: []fs.FileInfo{
				must(fs.Stat(testfs, "binary1.bin")),
				must(fs.Stat(testfs, "binary2.bin")),
				must(fs.Stat(testfs, "file1.txt")),
				must(fs.Stat(testfs, "file2.txt")),
				must(fs.Stat(testfs, "file3.txt")),
			},
			wantErr: false,
		},
		{
			name: "collect only binary files",
			args: args{
				fsys:  testfs,
				globs: []string{"*.bin"},
			},
			wantFiles: []fs.FileInfo{
				must(fs.Stat(testfs, "binary1.bin")),
				must(fs.Stat(testfs, "binary2.bin")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, err := collectFiles(tt.args.fsys, tt.args.globs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("collectFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, gotFiles, len(tt.wantFiles))
			assert.True(t, slices.EqualFunc(tt.wantFiles, gotFiles, func(a, b fs.FileInfo) bool {
				t.Logf("%s, %s => %v", a.Name(), b.Name(), a.Name() == b.Name())
				return a.Name() == b.Name()
			}))
		})
	}
}

func Test_collectDirs(t *testing.T) {
	type args struct {
		fsys fs.FS
	}
	tests := []struct {
		name    string
		args    args
		want    []fs.FileInfo
		wantErr bool
	}{
		{
			name: "collect all dirs",
			args: args{
				fsys: testfs,
			},
			want: []fs.FileInfo{
				must(fs.Stat(testfs, "dir1")),
				must(fs.Stat(testfs, "dir2")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := collectDirs(tt.args.fsys)
			if (err != nil) != tt.wantErr {
				t.Errorf("collectDirs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectDirs() = %v, want %v", got, tt.want)
			}
		})
	}
}
