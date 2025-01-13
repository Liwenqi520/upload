package upload

import "mime/multipart"

type Uploader interface {
	// UploadFile 上传文件
	UploadFile(file *multipart.FileHeader) (string, string, error)
	// DeleteFile 删除文件 p 文件位置路径 eg "uploads/temp/a.jpg"
	DeleteFile(p string) error
	// CopyFile 复制文件到指定文件夹位置
	// originalFile 原始文件位置 eg "uploads/temp/a.jpg"
	// targetFilePath 新文件夹位置 eg "user" 则会创建默认加上当前./uploads/
	CopyFile(originalFile string, targetFilePath string) (string, error)
}

func NewUploader(way string) Uploader {
	switch way {
	case "local":
		return &local{}
	default:
		return nil
	}
}
