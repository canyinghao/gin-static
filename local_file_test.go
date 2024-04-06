package static_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	static "github.com/soulteary/gin-static"
	"github.com/stretchr/testify/assert"
)

func TestLocalFile(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := os.CreateTemp(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("Gin Web Framework")
	f.Close()

	dir, filename := filepath.Split(f.Name())
	router := gin.New()
	router.Use(static.Serve("/", static.LocalFile(dir, true)))

	w := PerformRequest(router, "GET", "/"+filename)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Gin Web Framework")

	w = PerformRequest(router, "GET", "/")
	assert.Contains(t, w.Body.String(), `<a href="`+filename)

	w = PerformRequest(router, "GET", "/"+"../"+filename)
	assert.Equal(t, w.Code, 404)
	assert.Equal(t, w.Body.String(), "404 page not found")

	w = PerformRequest(router, "GET", "/"+"\\"+filename)
	assert.Equal(t, w.Code, 404)
	assert.Equal(t, w.Body.String(), "404 page not found")
}
