package users

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"
)

//func CreateTarGz(sourceDir, targetFile string) error {
//	target, err := os.Create(targetFile)
//	if err != nil {
//		return err
//	}
//	defer target.Close()
//
//	gw := gzip.NewWriter(target)
//	defer gw.Close()
//
//	tw := tar.NewWriter(gw)
//	defer tw.Close()
//	filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//		header, err := tar.FileInfoHeader(info, "")
//		if err != nil {
//			return err
//		}
//		header.Name = strings.TrimPrefix(strings.TrimPrefix(path, sourceDir), string(filepath.Separator))
//		header.Size = info.Size()
//		if err := tw.WriteHeader(header); err != nil {
//			return err
//		}
//		if info.IsDir() {
//			return nil
//		}
//
//		file, err := os.Open(path)
//		if err != nil {
//			return err
//		}
//		defer file.Close()
//		_, err = io.Copy(tw, file)
//		return err
//	})
//	return nil
//}

//	func CreateTarGz(sourceDir, targetFile string) error {
//		target, err := os.Create(targetFile)
//		if err != nil {
//			return err
//		}
//		defer target.Close()
//
//		gw := gzip.NewWriter(target)
//		defer gw.Close()
//
//		tw := tar.NewWriter(gw)
//
//		defer tw.Close()
//
//		err = filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
//			if err != nil {
//				return err
//			}
//
//			// Tạo header cho tệp
//			header, err := tar.FileInfoHeader(info, "")
//			if err != nil {
//				return err
//			}
//
//			// Đặt tên tệp trong tar
//			header.Name = strings.TrimPrefix(strings.TrimPrefix(path, sourceDir), string(filepath.Separator))
//			if err := tw.WriteHeader(header); err != nil {
//				return err
//			}
//			if info.IsDir() {
//				return nil
//			}
//
//			// Mở tệp
//			file, err := os.Open(path)
//			if err != nil {
//				return err
//			}
//			defer file.Close()
//
//			// Sử dụng bộ đệm để đọc tệp
//			buf := make([]byte, 1024*2)
//			for {
//				n, err := file.Read(buf)
//				if err != nil && err != io.EOF {
//					return err
//				}
//				if n == 0 {
//					break
//				}
//				if _, err := tw.Write(buf[:n]); err != nil {
//					log.Printf("failed to write %d bytes to tar for %s: %v", n, path, err)
//					return err
//				}
//			}
//			return nil
//		})
//
//		return err
//	}
func CreateTarGzVer2(sourceDir, targetFile string) error {
	absInputDir, err := filepath.Abs(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	// Tách tên thư mục cuối cùng khỏi đường dẫn
	dirName := filepath.Base(absInputDir)
	parentDir := filepath.Dir(absInputDir)
	cmd := exec.Command("tar", "-czvf", targetFile, "-C", parentDir, dirName)
	// Chạy câu lệnh và đợi kết quả
	err1 := cmd.Run()
	if err1 != nil {
		return fmt.Errorf("error creating tar.gz file: %w", err)
	}

	return nil
}
func GetFileTarGz(db *gorm.DB) func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		sourceDir := "./bundle"
		zipFileName := "bundle.tar.gz"

		//create file tar.gz
		err := CreateTarGzVer2(sourceDir, zipFileName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create file tar.gz",
			})
			return
		}
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Disposition", "attachment; filename="+zipFileName)
		ctx.Header("Content-Type", "application/gzip")
		ctx.File(zipFileName)

		//xoa file sau khi gui
		//os.Remove(zipFileName)
	}
}

//var mu sync.Mutex

func NotifyUpdate(db *gorm.DB) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println("Error creating watcher:", err)
			return
		}
		defer watcher.Close()

		err = watcher.Add("./bundle/data")
		if err != nil {
			log.Println("Error adding directory to watcher:", err)
			return
		}
		//go func() {

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					sourceDir := "./bundle"
					zipFileName := "bundle.tar.gz"

					//create file tar.gz
					err := CreateTarGzVer2(sourceDir, zipFileName)
					log.Println(err)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{
							"error": "Failed to create file tar.gz",
						})
						return
					}
					ctx.Header("Content-Description", "File Transfer")
					ctx.Header("Content-Disposition", "attachment; filename="+zipFileName)
					//ctx.Header("Content-Type", "application/gzip")
					ctx.Header("Content-Type", "application/vnd.openpolicyagent.bundles")
					ctx.File(zipFileName)
					return
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			case <-time.After(20 * time.Second):
				ctx.JSON(304, gin.H{
					"error": "timeout",
				})
				return
			}
		}
	}
}
