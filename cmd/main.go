package main

//go get github.com/gin-gonic/gin correr este comando antes de ejecutar el proyecto
// paraejecutar el programa go run ./cmd/main.go
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Ruta para servir el archivo HTML
	r.StaticFile("/", "./index.html") // Asegúrate de que index.html esté en la misma carpeta

	// Ruta para convertir con GET
	r.GET("/convert/:type", convertGet)

	// Ruta para convertir con POST
	r.POST("/convert", convert)

	// Iniciar el servidor en el puerto 8080
	r.Run(":8080")
}

// Función de conversión para GET
func convertGet(c *gin.Context) {
	typeConversion := c.Param("type")
	value := c.Query("value")

	if val, err := strconv.ParseFloat(value, 64); err == nil {
		result := performConversion(typeConversion, val)
		if result != nil {
			c.JSON(http.StatusOK, result)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de conversión no válido"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valor no válido"})
	}
}

// Función de conversión para POST
func convert(c *gin.Context) {
	var request struct {
		Type  string  `json:"type"`
		Value float64 `json:"value"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Solicitud inválida"})
		return
	}

	result := performConversion(request.Type, request.Value)
	if result != nil {
		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de conversión no válido"})
	}
}

// Función que realiza la conversión según el tipo
func performConversion(conversionType string, value float64) gin.H {
	conversions := map[string]func(float64) float64{
		"celsius-to-fahrenheit": func(v float64) float64 { return (v * 9 / 5) + 32 },
		"fahrenheit-to-celsius": func(v float64) float64 { return (v - 32) * 5 / 9 },
		"km-to-miles":           func(v float64) float64 { return v * 0.621371 },
		"miles-to-km":           func(v float64) float64 { return v / 0.621371 },
		"inches-to-cm":          func(v float64) float64 { return v * 2.54 },
		"cm-to-inches":          func(v float64) float64 { return v / 2.54 },
		"feet-to-meters":        func(v float64) float64 { return v * 0.3048 },
		"meters-to-feet":        func(v float64) float64 { return v / 0.3048 },
	}

	if convertFunc, exists := conversions[conversionType]; exists {
		result := convertFunc(value)
		return gin.H{"Type": conversionType, "Value": value, "Result": result}
	}
	return nil
}
