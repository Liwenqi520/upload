package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Liwenqi520/unique"
)

// 默认文件路径
var defaultUploadPath = "./public/uploads"

type local struct {
}

// removeTempFile 删除过期文件
func (l local) removeTempFile() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("remove temp file err", err)
		}
	}()
	diffTime := int64(3600)
	nowTime := time.Now().Unix()
	tempPath := path.Join(defaultUploadPath, "temp") //temp文件夹
	_ = filepath.Walk(tempPath, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fileTime := info.ModTime().Unix()
			if (nowTime - fileTime) > diffTime {
				os.Remove(p)
			}
		}
		return nil
	})
	return
}
func (l local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	//清除过期文件
	go l.removeTempFile()
	tempPath := path.Join(defaultUploadPath, "temp")
	ext := path.Ext(file.Filename)
	filename := unique.NewID() + ext
	mkdirErr := os.MkdirAll(tempPath, os.ModePerm)
	if mkdirErr != nil {
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	p := tempPath + "/" + filename
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()
	out, createErr := os.Create(p)
	if createErr != nil {
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close()
	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return p, filename, nil
}
func (l local) CopyFile(originalFile string, targetFilePath string) (string, error) {
	originalSinglePath := path.Join(originalFile)
	//判断是否存在
	originalFileInfo, err := os.Stat(originalSinglePath)
	if err != nil {
		return "", errors.New("原始文件不存在，function os.Stat() Filed, err:" + err.Error())
	}
	targetFilePath = path.Join(defaultUploadPath, targetFilePath)
	_, err = os.Stat(targetFilePath)
	if err != nil {
		fmt.Println("没有目标文件夹，正在创建...", targetFilePath)
		mkdirErr := os.MkdirAll(targetFilePath, os.ModePerm)
		if mkdirErr != nil {
			return "", errors.New("创建目标文件失败，function os.MkdirAll() Filed, err:" + err.Error())
		}
	}
	//newFile := targetFilePath + originalFileInfo.Name()
	newFileName := path.Join(targetFilePath, originalFileInfo.Name())
	oldFile, err := os.Open(originalSinglePath)
	//oldFile, err := os.Create(originalSinglePath) // 这个没有文件 只创建了
	if err != nil {
		return "", errors.New("create oldFile文件失败，function os.Create Filed, err:" + err.Error())
	}
	defer oldFile.Close()
	newFile, err := os.Create(newFileName)
	if err != nil {
		return "", errors.New("create NewFile失败，function os.Create Filed, err:" + err.Error())
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		return "", errors.New("复制文件失败，function io.Copy Filed, err:" + err.Error())
	}
	return newFileName, nil
}
func (l local) DeleteFile(p string) error {
	/*验证删除的位置*/
	if strings.Contains(p, strings.Trim(defaultUploadPath, "./")) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + p + err.Error())
		}
	}
	return nil
}
