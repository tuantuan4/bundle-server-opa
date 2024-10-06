package users

import (
	"archive/tar"
	"compress/gzip"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateTarGz(sourceDir, targetFile string) error {
	target, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer target.Close()

	gw := gzip.NewWriter(target)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()
	filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(strings.TrimPrefix(path, sourceDir), string(filepath.Separator))
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(tw, file)
		return err
	})
	return nil
}

func GetFileTarGz(db *gorm.DB) func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		sourceDir := "./bundle"
		zipFileName := "bundle.tar.gz"

		//create file tar.gz
		err := CreateTarGz(sourceDir, zipFileName)
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
					//select {
					//case changeOccurred <- true:
					//default:
					//}
					sourceDir := "./bundle"
					zipFileName := "bundle.tar.gz"

					//create file tar.gz
					err := CreateTarGz(sourceDir, zipFileName)
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
			case <-time.After(8 * time.Second):
				ctx.JSON(304, gin.H{
					"error": "timeout",
				})
				return
			}
		}
		//}()
		//	select {
		//	case <-done:
		//
		//	case <-time.After(5 * time.Second):
		//		// The goroutine is still running after 5 seconds
		//
		//		// Check if there was any change during the waiting period
		//		select {
		//		case <-changeOccurred:
		//			sourceDir := "./bundle"
		//			zipFileName := "bundle.tar.gz"
		//
		//			// Create tar.gz file.
		//			err := CreateTarGz(sourceDir, zipFileName)
		//			if err != nil {
		//				log.Println("Failed to create file tar.gz:", err)
		//				ctx.JSON(http.StatusInternalServerError, gin.H{
		//					"error": "Failed to create file tar.gz",
		//				})
		//				return
		//			}
		//
		//			// Return or serve the file
		//			ctx.Header("Content-Description", "File Transfer")
		//			ctx.Header("Content-Disposition", "attachment; filename="+zipFileName)
		//			ctx.Header("Content-Type", "application/gzip")
		//			ctx.File(zipFileName)
		//
		//			//os.Remove(zipFileName)
		//		default:
		//			ctx.JSON(304, gin.H{
		//				"message": "No change detected",
		//			})
		//		}
		//	}
	}
}
