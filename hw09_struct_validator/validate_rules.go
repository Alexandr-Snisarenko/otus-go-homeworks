package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrInvalidValidateRuleString = errors.New("validate rule string is invalid")
	ErrInvalidValidateRuleParam  = errors.New("validate param is incorrect")
	ErrUnknownValidateRule       = errors.New("unknown validate rule")
	ErrInvalidRuleDataType       = errors.New("rule from another data type")
	ErrUnknownDataType           = errors.New("unknown data type")
	// Validate Errors
	ErrStringLength   = errors.New("string length is not valid")
	ErrStringContent  = errors.New("string content is not valid")
	ErrStringNotInSet = errors.New("string is not in string set")
	ErrNumberValue    = errors.New("namber value is not valid")
	ErrNumberNotInSet = errors.New("number is not in number set")
	ErrDataValidate   = errors.New("data is not valid")
)

// правила валидации
type ValidateRule uint

const (
	vr_unknown ValidateRule = iota
	vr_len
	vr_regexp
	vr_in
	vr_min
	vr_max
)

// мапинг правил на названия в строке
var ValidateRuleNameMap = map[string]ValidateRule{
	"len":    vr_len,
	"regexp": vr_regexp,
	"in":     vr_in,
	"min":    vr_min,
	"max":    vr_max,
}

// перечень поддерживаемых правил для строк
var StringValidateRules = []ValidateRule{
	vr_len,    // "len:32" - длина строки должна быть ровно 32 символа;
	vr_regexp, // "regexp:\\d+" - согласно регулярному выражению строка должна состоять из цифр (\\ - экранирование слэша);
	vr_in,     //"in:foo,bar" - строка должна входить в множество строк {"foo", "bar"}
	vr_min,    // "min:10" - строка не может быть меньше 10 символов;
	vr_max,    // "max:20" - строка не может быть больше 20 символов;
}

// перечень поддерживаемых правил для чисел
var NumberValidateRules = []ValidateRule{
	vr_min, // "min:10" - число не может быть меньше 10;
	vr_max, // "max:20" - число не может быть больше 20;
	vr_in,  // "in:256,1024" - число должно входить в множество чисел {256, 1024};
}

// набор правил с параметрами кнтроля
type ValidateRuleSet map[ValidateRule]string

// конструктор для набора правил
func NewValidateRuleSet(validateString string) (error, *ValidateRuleSet) {
	vSet := make(ValidateRuleSet)
	// если строка правил не задана - возвращаем пустой RuleSet
	if validateString == "" {
		return nil, &vSet
	}

	// если строка не пустая - разбираем строку и формируем список правил валидации
	// делим строку на катерны по символу "|"
	for _, ruleStr := range strings.Split(validateString, "|") {
		// разделяем правило и контрольные значения по ":"
		sRule, sVal, ok := strings.Cut(ruleStr, ":")
		if !ok {
			return ErrInvalidValidateRuleString, nil
		}
		// если правила нет в общем списке правил - ошибка
		rule, ok := ValidateRuleNameMap[sRule]
		if !ok {
			return ErrUnknownValidateRule, nil
		}
		vSet[rule] = sVal
	}

	return nil, &vSet
}

// проверка данных по списку правил
func (v *ValidateRuleSet) CheckData(data interface{}) error {
	// получаем объект за интерфейсом
	rVal := reflect.ValueOf(data)

	// выполняем проверку данных в зависимости от типа данных
	switch rVal.Kind() {
	case reflect.String:
		return checkStringData(v, rVal.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return checkNumberData(v, rVal.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return checkNumberData(v, rVal.Uint())
	default:
		return ErrUnknownDataType
	}
}

// методы для провеки значений по типам реализуются в виде отдельных функций
// ввиду ограничения для дженериков в методах
func checkStringData(rSet *ValidateRuleSet, data string) error {
	for rule, val := range *rSet {
		if !slices.Contains(StringValidateRules, rule) {
			return ErrInvalidRuleDataType
		}

		switch rule {
		case vr_len:
			chkVal, err := strconv.Atoi(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if len(data) != chkVal {
				return fmt.Errorf("%w not match given size", ErrStringLength)
			}
		case vr_min:
			chkVal, err := strconv.Atoi(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if len(data) < chkVal {
				return fmt.Errorf("%w  less then 'min' size", ErrStringLength)
			}
		case vr_max:
			chkVal, err := strconv.Atoi(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if len(data) > chkVal {
				return fmt.Errorf("%w greate then 'max' size", ErrStringLength)
			}
		case vr_regexp:
			rexp, err := regexp.Compile(val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}

			if !rexp.MatchString(data) {
				return ErrStringContent
			}
		case vr_in:
			if !slices.Contains(strings.Split(val, ","), data) {
				return ErrStringNotInSet
			}
		}
	}

	return nil
}

// // generics ////
// функция парсировки данных из строки в number (int64 или uint64)
func parseAnyInt[T int64 | uint64](val string) (T, error) {
	var (
		parsVal T
		i64     int64
		ui64    uint64
		err     error
	)
	// определяем тип данных от T
	dataType := fmt.Sprintf("%T", parsVal)

	switch dataType {
	case "int64":
		i64, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return parsVal, ErrInvalidValidateRuleParam
		}
		parsVal = T(i64)
	case "uint64":
		ui64, err = strconv.ParseUint(val, 10, 64)
		if err != nil {
			return parsVal, ErrInvalidValidateRuleParam
		}
		parsVal = T(ui64)
	default:
		return parsVal, ErrUnknownDataType
	}

	return parsVal, nil
}

// функця проверки параметров типа int64 или uint64
func checkNumberData[T int64 | uint64](rSet *ValidateRuleSet, data T) error {
	for rule, val := range *rSet {
		if !slices.Contains(NumberValidateRules, rule) {
			return ErrInvalidRuleDataType
		}
		switch rule {
		case vr_min:
			chkVal, err := parseAnyInt[T](val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}
			if data < chkVal {
				return fmt.Errorf("%w less then 'min' value", ErrNumberValue)
			}
		case vr_max:
			chkVal, err := parseAnyInt[T](val)
			if err != nil {
				return ErrInvalidValidateRuleParam
			}

			if data > chkVal {
				return fmt.Errorf("%w greate then 'max' value", ErrNumberValue)
			}
		case vr_in:
			// разбираем строку параметров на элементы
			sSlice := strings.Split(val, ",")
			// создаем пустой слайс под параметры типа int64 нужного размера
			chkSlc := make([]T, 0, len(sSlice))
			// парсим слайс строковых параметров и переводим его в слайс int64
			for _, s := range sSlice {
				chkVal, err := parseAnyInt[T](s)
				if err != nil {
					return ErrInvalidValidateRuleParam
				}
				chkSlc = append(chkSlc, chkVal)
			}
			// проверяем правило "in"
			if !slices.Contains(chkSlc, data) {
				return ErrNumberNotInSet
			}
		}
	}
	return nil
}
