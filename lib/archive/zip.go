package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/rip0532/mfano/lib/constant"
	logger "github.com/rip0532/mfano/lib/log"
)

type Zip struct {
}

func NewZip() *Zip {
	return &Zip{}
}

func (z *Zip) UnZip(dst, filename string) (err error) {
	zr, err := zip.OpenReader(constant.HomeDir + "/" + filename)
	defer zr.Close()
	if err != nil {
		return
	}

	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			log.Println(err.Error())
			return err
		}
	}

	for _, file := range zr.File {
		var filename string
		if !utf8.ValidString(file.Name) {
			i := bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			filename = string(content)
		} else {
			filename = file.Name
		}
		// 处理macos压缩元数据
		if strings.Contains(filename, "__MACOSX") {
			continue
		}
		path := filepath.Join(dst, filename)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, 0755); err != nil {
				logger.Error.Println(path)
				return err
			}
			continue
		}

		fr, err := file.Open()
		if err != nil {
			return err
		}

		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
		fmt.Printf("unzip file:%s successes！\n", path)
		fw.Close()
		fr.Close()
	}
	return nil
}
