package validate

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strconv"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/vi"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var trans ut.Translator

type RequestError struct {
	Fields []Field `json:"fields"`
}

type Field struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Errs    string `json:"errs"`
}

func (r RequestError) Error() string {
	payload, _ := json.Marshal(r)
	return string(payload)
}

var translations = map[string]map[string]string{
	"pt": {
		"required": "{0} é um campo obrigatório.",
		"email":    "{0} inválido.",
		"gte":      "{0} deve ter {1} ou mais caracteres.",
		"lte":      "{0} deve ter {1} ou menos caracteres.",
	},
}

var dicts = map[string]map[string]string{
	"pt": {
		"Title":       "Título",
		"Description": "Descrição",
	},
}

func translationFunc(t ut.Translator, fe validator.FieldError) string {
	field, err := t.T(fe.Field())
	if err != nil {
		field = fe.Field()
	}
	msg, err := t.T(fe.Tag(), field, fe.Param())
	if err != nil {
		return fe.Error()
	}
	return msg
}

func InitializeValidator() *validator.Validate {
	validate := validator.New()
	enLocale := en.New()
	utrans := ut.New(enLocale, enLocale, vi.New())
	trans, _ = utrans.GetTranslator("pt")

	validate.RegisterValidation("ISO8601date", isISO8601Date)
	validate.RegisterValidation("PhoneFormat", phoneFormat)
	validate.RegisterValidation("CPForCNPJ", validateCPForCNPJ)

	for locale, dict := range dicts {
		engine, _ := utrans.FindTranslator(locale)
		for key, trans := range dict {
			_ = engine.Add(key, trans, false)
		}
	}

	// for locale, translation := range translations {
	for locale, translation := range translations {
		engine, _ := utrans.FindTranslator(locale)
		for tag, trans := range translation {
			_ = validate.RegisterTranslation(tag, engine, func(t ut.Translator) error {
				return t.Add(tag, trans, false)
			}, translationFunc)
		}
	}

	// pt_translations.RegisterDefaultTranslations(validate, trans)
	return validate
}

func validateCPForCNPJ(fl validator.FieldLevel) bool {
	// Remove non-digit characters from the input
	regex := regexp.MustCompile(`\D`)
	number := regex.ReplaceAllString(fl.Field().String(), "")

	if len(number) == 11 {
		return validateCPF(fl.Field().String())
	} else if len(number) == 14 {
		return validateCNPJ(number)
	}

	return true
}

func validateCPF(cpf string) bool {
	// Remove non-digit characters from the input
	regex := regexp.MustCompile(`\D`)
	cpf = regex.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	digit1, _ := strconv.Atoi(string(cpf[9]))
	digit2, _ := strconv.Atoi(string(cpf[10]))

	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}

	remainder := sum % 11
	if remainder < 2 && digit1 != 0 || remainder >= 2 && digit1 != 11-remainder {
		return false
	}

	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}

	remainder = sum % 11
	if remainder < 2 && digit2 != 0 || remainder >= 2 && digit2 != 11-remainder {
		return false
	}

	return true
}

func validateCNPJ(cnpj string) bool {
	// Remove non-digit characters from the input
	regex := regexp.MustCompile(`\D`)
	cnpj = regex.ReplaceAllString(cnpj, "")

	if len(cnpj) != 14 {
		return false
	}

	var weights = [12]int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights[i]
	}

	remainder := sum % 11
	digit1, _ := strconv.Atoi(string(cnpj[12]))
	if remainder < 2 && digit1 != 0 || remainder >= 2 && digit1 != 11-remainder {
		return false
	}

	sum = 0
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights[i]
	}

	remainder = sum % 11
	digit2, _ := strconv.Atoi(string(cnpj[13]))
	if remainder < 2 && digit2 != 0 || remainder >= 2 && digit2 != 11-remainder {
		return false
	}

	return true
}

func isISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
	return ISO8601DateRegex.MatchString(fl.Field().String())
}

func phoneFormat(fl validator.FieldLevel) bool {
	PhoneFormatRegexString := "^\\(\\d{2}\\) \\d{5}-\\d{4}$"
	PhoneFormatRegex := regexp.MustCompile(PhoneFormatRegexString)
	return PhoneFormatRegex.MatchString(fl.Field().String())
}

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}

		return nil
	}
}

func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}

	return msg
}

func HandleValidatorFieldError(data interface{}, err error) error {
	errs := err.(validator.ValidationErrors)

	requestError := RequestError{}

	for _, e := range errs {
		field := Field{
			Field:   getFormatedField(data, e.Field()),
			Message: e.Translate(trans),
			Errs:    e.Error(),
		}

		requestError.Fields = append(requestError.Fields, field)
	}

	return requestError
}

func getFormatedField(data interface{}, field string) string {
	if field, ok := reflect.TypeOf(data).Elem().FieldByName(field); ok {
		return field.Tag.Get("json")
	}

	return ""
}
