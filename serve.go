package static

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

func ServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return Serve(urlPrefix, LocalFile(root, false))
}

// Serve returns a middleware handler that serves static files in the given directory.
func Serve(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	return ServeCached(urlPrefix, fs, -1)
}

// ServeCached returns a middleware handler that similar as Serve
// but with the Cache-Control Header set as passed in the cacheAge parameter
func ServeCached(urlPrefix string, fs ServeFileSystem, cacheAge int) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {

			if cacheAge > -1 {
				fl, err := fs.Open(c.Request.URL.Path)
				if err == nil {
					defer fl.Close()
					// 生成并设置 ETag 头
					eTag := generateETag(fl)
					c.Writer.Header().Add("ETag", eTag)

					// 检查 If-None-Match 头与生成的 ETag 是否匹配，若匹配则返回 304 Not Modified
					if match := c.GetHeader("If-None-Match"); match != "" {
						if match == eTag {
							log.Printf("Cache hit for: %s", c.Request.URL.Path)
							c.Status(http.StatusNotModified)
							c.Abort()
							return
						}
					}
				}

				c.Writer.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", cacheAge))

			}
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

// generateETag 根据文件内容生成一个 ETag
func generateETag(f http.File) string {

	fileContent, err := ReadFile(f)
	if err != nil {
		log.Printf("Error reading file %v", err)
		return ""
	}

	h := sha1.New()
	h.Write(fileContent)
	return fmt.Sprintf("\"%s\"", hex.EncodeToString(h.Sum(nil)))
}

func ReadFile(f http.File) ([]byte, error) {

	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++ // one byte for final read at EOF
	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}

		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
	}
}
