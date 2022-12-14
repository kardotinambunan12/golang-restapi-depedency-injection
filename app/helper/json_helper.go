package helper

import (
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/spf13/viper"
)

// func ReadFromRequestBody(ctx *fiber.Ctx, result interface{}) {
// 	decoder := json.NewDecoder(ctx.Body)
// 	err := decoder.Decode(result)
// 	errorhandler.PanicIfNeeded(err)
// }

// StorageDirectory func
func StorageDirectory() string {
	storage := viper.GetString("STORAGE_DIRECTORY")
	if storage == "" {
		currentDir, _ := os.Getwd()
		storage = currentDir + "/uploads"
	}

	storage = strings.ReplaceAll(storage, "//", "/")

	if !DirExists(storage) && viper.GetBool("STORAGE_CREATE") {
		os.MkdirAll(storage+"/resize", 0755)
	}

	return storage
}

// DirExists func
func DirExists(dirname string) bool {
	info, err := os.Stat(filepath.FromSlash(dirname))
	return nil == err && !os.IsNotExist(err) && info.IsDir()
}

// FileExists func
func FileExists(filename string) bool {
	info, err := os.Stat(filepath.FromSlash(filename))
	return nil == err && !os.IsNotExist(err) && !info.IsDir()
}

// GetMimeFile func
func GetMimeFile(filepath string) types.Type {
	var ft types.Type = types.Unknown
	if FileExists(filepath) {
		data, err := ioutil.ReadFile(filepath)
		if nil == err {
			typ, err := filetype.Get(data)
			if nil == err {
				ft = typ
			}
		}
		if ft == types.Unknown {
			ctype := http.DetectContentType(data)
			ctypes := strings.Split(ctype, ";")
			if len(ctypes) > 0 {
				ft.MIME = types.MIME{
					Type:    ctypes[0],
					Subtype: ctypes[0],
					Value:   ctypes[0],
				}
				ft.MIME.Type = ctypes[0]
				ft.MIME.Value = ctypes[0]
			}

		}
	}
	names := strings.Split(filepath, ".")
	ln := len(names)
	if ln > 1 {
		exts := []string{}
		if ln >= 2 && strings.ToLower(names[ln-2]) == "tar" {
			exts = append(exts, "tar")
		}
		exts = append(exts, strings.ToLower(names[ln-1]))
		ft.Extension = strings.Join(exts, ".")
	}

	if ft == types.Unknown {
		if typ, err := filetype.MatchFile(filepath); nil != err {
			ft = typ
		}
	}

	return ft
}

// GetImageScaleSize func
func GetImageScaleSize(filename string) (w int, h int, err error) {
	file, _ := os.Open(filename)
	img, err := jpeg.DecodeConfig(file)
	file.Close()
	if nil != err {
		file, _ = os.Open(filename)
		img, err = png.DecodeConfig(file)
		file.Close()
		if nil != err {
			return 0, 0, err
		}
	}

	return img.Width, img.Height, nil
}
