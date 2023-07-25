package postgres

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/lib/pq"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

// Constraint is a custom database check constraint you've defined, like "CHECK
// balance > 0". Postgres doesn't define a very useful message for constraint
// failures (new row for relation "accounts" violates check constraint), so you
// can define your own. The Name should be the name of the constraint in the
// database. Define GetError to provide your own custom error handler for this
// constraint failure, with a custom message.
type Constraint struct {
	Name     string
	GetError func(*pq.Error) *PgxError
}

var constraintMap = map[string]*Constraint{}
var constraintMu sync.RWMutex

// RegisterConstraint tells dberror about your custom constraint and its error
// handling. RegisterConstraint panics if you attempt to register two
// constraints with the same name.
func RegisterConstraint(c *Constraint) {
	constraintMu.Lock()
	defer constraintMu.Unlock()
	if _, dup := constraintMap[c.Name]; dup {
		panic("dberror: RegisterConstraint called twice for name " + c.Name)
	}
	constraintMap[c.Name] = c
}

// capitalize the first letter in the string
func capitalize(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return fmt.Sprintf("%c", unicode.ToTitle(r)) + s[size:]
}

var columnFinder = regexp.MustCompile(`Key \((.+)\)=`)
var valueFinder = regexp.MustCompile(`Key \(.+\)=\((.+)\)`)

// findColumn finds the column in the given pq Detail error string. If the
// column does not exist, the empty string is returned.
//
// detail can look like this:
//
//	Key (id)=(3c7d2b4a-3fc8-4782-a518-4ce9efef51e7) already exists.
func findColumn(detail string) string {
	results := columnFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	} else {
		return results[1]
	}
}

// findColumn finds the column in the given pq Detail error string. If the
// column does not exist, the empty string is returned.
//
// detail can look like this:
//
//	Key (id)=(3c7d2b4a-3fc8-4782-a518-4ce9efef51e7) already exists.
func findValue(detail string) string {
	results := valueFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	} else {
		return results[1]
	}
}

var foreignKeyFinder = regexp.MustCompile(`not present in table "(.+)"`)

// findForeignKeyTable finds the referenced table in the given pq Detail error
// string. If we can't find the table, we return the empty string.
//
// detail can look like this:
//
//	Key (account_id)=(91f47e99-d616-4d8c-9c02-cbd13bceac60) is not present in table "accounts"
func findForeignKeyTable(detail string) string {
	results := foreignKeyFinder.FindStringSubmatch(detail)
	if len(results) < 2 {
		return ""
	}
	return results[1]
}

var parentTableFinder = regexp.MustCompile(`update or delete on table "([^"]+)"`)

func findParentTable(message string) string {
	match := parentTableFinder.FindStringSubmatch(message)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

type PgxError struct {
	Message    string `json:"message"`
	Code       string `json:"code"`
	Constraint string `json:"constraint"`
	Severity   string `json:"severity"`
	Routine    string `json:"routine"`
	Table      string `json:"table"`
	Detail     string `json:"detail"`
	Column     string `json:"column"`
}

func (dbe *PgxError) Error() string {
	return dbe.Message
}

func InitializeTranslations() {
	bundle = i18n.NewBundle(language.BrazilianPortuguese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err := bundle.LoadMessageFile("pt_BR.json")
	fmt.Println(err)
}

func GetError(err error) error {
	switch sqlxerr := err.(type) {
	case *pq.Error:
		tags := map[string]interface{}{
			"Message": sqlxerr.Message,
		}

		localizer := i18n.NewLocalizer(bundle, "pt_BR") // Set the desired locale

		translatedError, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{ //
				ID:    string(sqlxerr.Code),
				Other: sqlxerr.Message,
			},
			TemplateData: tags,
		})

		return &PgxError{
			Message:    fmt.Sprintf(translatedError, returnField(string(sqlxerr.Code), sqlxerr)),
			Code:       string(sqlxerr.Code),
			Column:     sqlxerr.Column,
			Detail:     sqlxerr.Message,
			Constraint: sqlxerr.Constraint,
			Table:      sqlxerr.Table,
			Routine:    sqlxerr.Routine,
			Severity:   sqlxerr.Severity,
		}
	default:
		return err
	}
}

func returnField(code string, err *pq.Error) string {
	switch code {
	case "22003":
		return findColumn(err.Detail)
	case "22P02":
		return findValue(err.Detail)
	case "23502":
		re := regexp.MustCompile(`column "(.+?)" of relation "(.+?)"`)
		match := re.FindStringSubmatch(err.Message)

		if len(match) > 1 {
			columnName := match[1]
			return columnName
		}
		break
	case "23503":
		re := regexp.MustCompile(`table "(.+?)"`)
		match := re.FindStringSubmatch(err.Detail)

		if len(match) > 1 {
			columnName := match[1]
			return columnName
		}
		break
	case "23505":
		re := regexp.MustCompile(`constraint "(.+?)"`)
		match := re.FindStringSubmatch(err.Message)

		if len(match) > 1 {
			columnName := match[1]
			return columnName
		}
		break
	case "23514":
		return findParentTable(err.Detail)
	case "42P01":
		return findParentTable(err.Detail)
	case "42P02":
		return findColumn(err.Detail)
	case "42703":
		re := regexp.MustCompile(`column "(.+?)" of relation "(.+?)"`)
		match := re.FindStringSubmatch(err.Message)

		if len(match) > 1 {
			columnName := match[1]
			return columnName
		}
		break
	case "42702":
		return findColumn(err.Detail)
	}

	return ""
}
