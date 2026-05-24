package validator

import (
    "reflect"
    "regexp"
    "strings"
    "time"

    "github.com/go-playground/validator/v10"
)

type CustomValidator struct {
    Validate *validator.Validate
}

func NewValidator() *CustomValidator {
    v := validator.New()
    
    // Register custom validation functions
    v.RegisterValidation("phone", validatePhone)
    v.RegisterValidation("password", validatePassword)
    v.RegisterValidation("date", validateDate)
    v.RegisterValidation("future_date", validateFutureDate)
    v.RegisterValidation("past_date", validatePastDate)
    v.RegisterValidation("timezone", validateTimezone)
    v.RegisterValidation("duration", validateDuration)
    v.RegisterValidation("tags", validateTags)
    v.RegisterValidation("url_or_empty", validateURLOrEmpty)
    
    v.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })
    
    return &CustomValidator{Validate: v}
}

func (cv *CustomValidator) ValidateStruct(s interface{}) map[string]string {
    errors := make(map[string]string)
    
    err := cv.Validate.Struct(s)
    if err == nil {
        return nil
    }
    
    for _, err := range err.(validator.ValidationErrors) {
        field := err.Field()
        tag := err.Tag()
        param := err.Param()
        
        message := getErrorMessage(field, tag, param)
        errors[strings.ToLower(field)] = message
    }
    
    return errors
}

func (cv *CustomValidator) ValidateVar(field interface{}, tag string) error {
    return cv.Validate.Var(field, tag)
}

func getErrorMessage(field, tag, param string) string {
    switch tag {
    case "required":
        return field + " is required"
    case "email":
        return field + " must be a valid email address"
    case "min":
        return field + " must be at least " + param + " characters long"
    case "max":
        return field + " must be at most " + param + " characters long"
    case "len":
        return field + " must be exactly " + param + " characters long"
    case "gte":
        return field + " must be greater than or equal to " + param
    case "lte":
        return field + " must be less than or equal to " + param
    case "phone":
        return field + " must be a valid phone number"
    case "password":
        return field + " must contain at least 8 characters, one uppercase, one lowercase, one number and one special character"
    case "date":
        return field + " must be a valid date (YYYY-MM-DD)"
    case "future_date":
        return field + " must be a future date"
    case "past_date":
        return field + " must be a past date"
    case "timezone":
        return field + " must be a valid timezone"
    case "duration":
        return field + " must be a positive duration"
    case "tags":
        return field + " must be a comma-separated list"
    case "oneof":
        return field + " must be one of [" + param + "]"
    case "numeric":
        return field + " must be numeric"
    case "alpha":
        return field + " must contain only letters"
    case "alphanum":
        return field + " must contain only letters and numbers"
    case "url":
        return field + " must be a valid URL"
    case "url_or_empty":
        return field + " must be a valid URL or empty"
    case "uuid":
        return field + " must be a valid UUID"
    default:
        return field + " is invalid"
    }
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    if phone == "" {
        return true 
    }
    
    matched, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, phone)
    return matched
}

func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    if password == "" {
        return false
    }
    
    if len(password) < 8 {
        return false
    }
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
    
    return hasUpper && hasLower && hasNumber && hasSpecial
}

func validateDate(fl validator.FieldLevel) bool {
    dateStr := fl.Field().String()
    if dateStr == "" {
        return true
    }
    
    _, err := time.Parse("2006-01-02", dateStr)
    return err == nil
}

func validateFutureDate(fl validator.FieldLevel) bool {
    dateStr := fl.Field().String()
    if dateStr == "" {
        return true
    }
    
    t, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return false
    }
    
    return t.After(time.Now())
}

func validatePastDate(fl validator.FieldLevel) bool {
    dateStr := fl.Field().String()
    if dateStr == "" {
        return true
    }
    
    t, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        return false
    }
    
    return t.Before(time.Now())
}

func validateTimezone(fl validator.FieldLevel) bool {
    tz := fl.Field().String()
    if tz == "" {
        return true
    }
    
    _, err := time.LoadLocation(tz)
    return err == nil
}

func validateDuration(fl validator.FieldLevel) bool {
    duration := fl.Field().Int()
    return duration > 0
}

func validateTags(fl validator.FieldLevel) bool {
    tags := fl.Field().String()
    if tags == "" {
        return true
    }
    
    parts := strings.Split(tags, ",")
    for _, part := range parts {
        part = strings.TrimSpace(part)
        if part == "" {
            return false
        }
        matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, part)
        if !matched {
            return false
        }
    }
    return true
}

func validateURLOrEmpty(fl validator.FieldLevel) bool {
    url := fl.Field().String()
    if url == "" {
        return true
    }
    
    matched, _ := regexp.MatchString(`^https?://[a-zA-Z0-9][a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`, url)
    return matched
}

var ValidationRules = struct {
    Email    string
    Password string
    Name     string
    Phone    string
    Bio      string
    UUID     string
}{
    Email:    "required,email,max=255",
    Password: "required,password",
    Name:     "required,min=2,max=100,alpha",
    Phone:    "omitempty,phone",
    Bio:      "omitempty,max=500",
    UUID:     "required,uuid",
}

func ValidateEmail(email string) bool {
    v := NewValidator()
    return v.ValidateVar(email, "required,email") == nil
}

func ValidatePassword(password string) bool {
    v := NewValidator()
    return v.ValidateVar(password, "password") == nil
}

func SanitizeInput(input string) string {
    input = strings.TrimSpace(input)
    
    input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")
    
    input = strings.ReplaceAll(input, "&", "&amp;")
    input = strings.ReplaceAll(input, "<", "&lt;")
    input = strings.ReplaceAll(input, ">", "&gt;")
    input = strings.ReplaceAll(input, "\"", "&quot;")
    input = strings.ReplaceAll(input, "'", "&#39;")
    
    return input
}

type ValidationResult struct {
    Valid  bool
    Errors map[string]string
}

func ValidateRequest(req interface{}) *ValidationResult {
    v := NewValidator()
    errors := v.ValidateStruct(req)
    
    return &ValidationResult{
        Valid:  len(errors) == 0,
        Errors: errors,
    }
}