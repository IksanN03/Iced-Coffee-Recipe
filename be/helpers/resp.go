package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type MessageType struct {
	Success string `json:"success,omitempty"`
	Warning string `json:"warning,omitempty"`
	Danger  string `json:"danger,omitempty"`
}

type APIResponseNew struct {
	Message MessageType       `json:"message"`
	Data    interface{}       `json:"data"`
	Error   map[string]string `json:"error"`
}

func getStatusFromError(err error) int {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound
	case errors.As(err, &validator.ValidationErrors{}):
		return http.StatusBadRequest
	case strings.Contains(err.Error(), "Bind"):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			errors[field] = fmt.Sprintf("Field validation for '%s' failed on the '%s'", field, e.Tag())
		}
		return errors
	}

	return nil
}

func NewAPIResponse(c *gin.Context, data interface{}, err error, attribute_or_entity string, status_code int, message string) {
	response := APIResponseNew{}
	default_attribute_or_entity := attribute_or_entity
	method := c.Request.Method
	if status_code == 0 {
		status_code = http.StatusOK
		if err != nil {
			status_code = getStatusFromError(err)
		}
	}

	// Get the type name and check if it's a slice
	var struct_name string
	if data != nil {
		t := reflect.TypeOf(data)
		isSlice := false
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Slice {
			isSlice = true
			t = t.Elem()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
		}

		name := t.Name()
		if strings.Contains(name, "") {
			// Multi word - convert to camelCase
			words := strings.Split(name, "")
			for i := 0; i < len(words); i++ {
				if i == 0 {
					words[i] = strings.ToLower(words[i])
				}
			}
			struct_name = strings.Join(words, "")
		} else {
			// Single word - all lowercase
			struct_name = strings.ToLower(name)
		}

		// Add 's' for slices
		if isSlice {
			struct_name += "s"
		}
	}

	if err == nil {
		if struct_name != "h" {
			attribute_or_entity = struct_name
			if attribute_or_entity != "" {
				struct_name = attribute_or_entity
			}
		}
	}
	if err != nil {
		if struct_name != "" {
			attribute_or_entity = struct_name
		}
		if validationErrors := formatValidationError(err); validationErrors != nil {
			response.Error = validationErrors
		} else {
			response.Error = map[string]string{
				attribute_or_entity: err.Error(),
			}
		}
	}
	if default_attribute_or_entity != "" {
		attribute_or_entity = default_attribute_or_entity
	}

	msg := MessageType{}
	switch status_code {
	case http.StatusOK:
		switch method {
		case "GET":
			if message == "" {
				msg.Success = fmt.Sprintf("Berhasil mengambil data %s", attribute_or_entity)
			} else {
				msg.Success = message
			}
		case "POST":
			if message == "" {
				msg.Success = fmt.Sprintf("Berhasil menambah data %s", attribute_or_entity)
			} else {
				msg.Success = message
			}
		case "PUT":
			if message == "" {
				msg.Success = fmt.Sprintf("Berhasil mengubah data %s", attribute_or_entity)
			} else {
				msg.Success = message
			}
		case "DELETE":
			if message == "" {
				msg.Success = fmt.Sprintf("Berhasil menghapus data %s", attribute_or_entity)
			} else {
				msg.Success = message
			}
		}
	case http.StatusNoContent:
		msg.Success = fmt.Sprintf("Data %s belum tersedia", attribute_or_entity)
	case http.StatusBadRequest:
		if message == "" {
			msg.Warning = "Permintaan tidak valid"
		} else {
			msg.Warning = message
		}
	case http.StatusUnauthorized:
		if message == "" {
			msg.Danger = "Tidak memiliki akses"
		} else {
			msg.Danger = message
		}
	case http.StatusNotFound:
		if message == "" {
			msg.Warning = fmt.Sprintf("Data %s tidak ditemukan", attribute_or_entity)
		} else {
			msg.Warning = message
		}
	case http.StatusPreconditionRequired:
		msg.Warning = message
	case http.StatusServiceUnavailable:
		msg.Danger = message
	case http.StatusNotAcceptable:
		msg.Danger = message
	case http.StatusForbidden:
		msg.Danger = message
	case http.StatusBadGateway:
		msg.Danger = "Server tidak dapat merespons permintaan, mohon tunggu beberapa saat"
	case http.StatusInternalServerError:
		if message == "" {
			msg.Danger = "Terjadi kesalahan pada server"
		} else {
			msg.Danger = message
		}
	}
	response.Message = msg

	if data != nil && err == nil {
		fmt.Println("ss")
		// Wrap the data in a map with plural struct name as key
		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Slice && v.Len() == 0 {
			data = nil
		}
		if struct_name == "h" {
			response.Data = data
		} else {

			if default_attribute_or_entity != "" {
				struct_name = default_attribute_or_entity
			}

			response.Data = map[string]interface{}{
				struct_name: data,
			}
		}

	}

	c.JSON(status_code, response)
}
