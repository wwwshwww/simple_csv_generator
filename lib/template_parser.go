package lib

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/samber/mo"
	"gopkg.in/yaml.v3"
)

type ColumnType string

const (
	ColumnTypeInt                  ColumnType = "INT"
	ColumnTypeFloat                ColumnType = "FLOAT"
	ColumnTypeBool                 ColumnType = "BOOL"
	ColumnTypeDatetime             ColumnType = "DATETIME"
	ColumnTypeString               ColumnType = "STRING"
	ColumnTypeMultilineString      ColumnType = "MULTILINE_STRING"
	ColumnTypeURL                  ColumnType = "URL"
	ColumnTypeArrayInt             ColumnType = "ARRAY_INT"
	ColumnTypeArrayFloat           ColumnType = "ARRAY_FLOAT"
	ColumnTypeArrayBool            ColumnType = "ARRAY_BOOL"
	ColumnTypeArrayDatetime        ColumnType = "ARRAY_DATETIME"
	ColumnTypeArrayString          ColumnType = "ARRAY_STRING"
	ColumnTypeArrayMultilineString ColumnType = "ARRAY_MULTILINE_STRING"
	ColumnTypeArrayURL             ColumnType = "ARRAY_URL"
)

func (s *ColumnType) UnmarshalYAML(value *yaml.Node) error {
	var tml string
	if err := value.Decode(&tml); err != nil {
		return err
	}

	switch ColumnType(tml) {
	case ColumnTypeInt, ColumnTypeFloat, ColumnTypeBool, ColumnTypeDatetime, ColumnTypeString, ColumnTypeMultilineString, ColumnTypeURL, ColumnTypeArrayInt, ColumnTypeArrayFloat, ColumnTypeArrayBool, ColumnTypeArrayDatetime, ColumnTypeArrayString, ColumnTypeArrayMultilineString, ColumnTypeArrayURL:
		*s = ColumnType(tml)
		return nil
	default:
		return fmt.Errorf("unknown type: %v", tml)
	}

}

type Config struct {
	Columns []Column `yaml:"columns"`
}

type Column struct {
	Name    string      `yaml:"name"`
	Type    ColumnType  `yaml:"type"`
	Choices ChoiceField `yaml:"choices,omitempty"`
}

type ChoiceField struct {
	values mo.Either[[]string, [][]string]
}

func (c *ChoiceField) UnmarshalYAML(value *yaml.Node) error {
	var single []string
	var multi [][]string
	if err := value.Decode(&single); err == nil {
		c.values = mo.Left[[]string, [][]string](single)
		return nil
	}
	if err := value.Decode(&multi); err == nil {
		c.values = mo.Right[[]string, [][]string](multi)
		return nil
	}
	return errors.New("invalid format for choices")
}

func ParseColumns(yaml_path string) mo.Result[[]Column] {
	file, err := os.Open(yaml_path)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		mo.Errf[[]Column]("failed to decode: %v", err)
	}

	return mo.Ok(config.Columns)
}

func bulkUnmarshal[T1, T2 any](target []T1, convFn func(T1) (T2, error)) mo.Result[[]T2] {
	errs := []error{}
	converted := []T2{}
	for _, v := range target {
		if c, err := convFn(v); err != nil {
			errs = append(errs, err)
		} else {
			converted = append(converted, c)
		}
	}
	if len(errs) > 0 {
		return mo.Errf[[]T2]("failed to convert: %+v", errs)
	}
	return mo.Ok(converted)
}

func UnmarshalIntChoices(target []string) mo.Result[[]int] {
	return bulkUnmarshal(target, strconv.Atoi)
}

func UnmarshalFloatChoices(target []string) mo.Result[[]float64] {
	return bulkUnmarshal(target, func(v string) (float64, error) {
		return strconv.ParseFloat(v, 64)
	})
}

func UnmarshalBoolChoices(target []string) mo.Result[[]bool] {
	return bulkUnmarshal(target, strconv.ParseBool)
}

func UnmarshalDatetimeChoices(target []string) mo.Result[[]time.Time] {
	return bulkUnmarshal(target, func(v string) (time.Time, error) {
		return time.Parse(time.RFC3339, v)
	})
}
