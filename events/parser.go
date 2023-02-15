package events

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const expectedStructTag = "rcon"

var typeCache = map[string]RCONType{}
var typeCacheMu = sync.RWMutex{}

type Columns []string

func (s Columns) GenerateDefaultHeader() string {
	var sb strings.Builder

	for i, col := range s {
		sb.WriteString(col)
		if isLast := len(s)-1 == i; !isLast {
			sb.WriteRune(Delimiter)
		}
		// New line isn't needed, we'll do that elsewhere
	}

	return sb.String()
}

type RCONType struct {
	Name        string
	Columns     Columns
	ColumnTypes map[string]reflect.Type

	columnLookup map[string]struct{}
}

func (s RCONType) Parse(header, input string) (map[string]any, error) {
	out := map[string]any{}

	// Split our columns out so we know the order
	headerParts := strings.Fields(header)
	if len(headerParts) > len(s.Columns) {
		return nil, fmt.Errorf("events/RCONType: too many parts in event header. Expected < %d got %d", len(s.Columns), len(headerParts))
	}

	// Sanity checking on the header parts
	for _, part := range headerParts {
		if _, ok := s.columnLookup[part]; !ok {
			return nil, fmt.Errorf("events/RCONType: unexpected column in header %s", part)
		}
	}

	parts := strings.Split(input, string(Delimiter))
	if len(parts) > len(s.Columns) {
		return nil, fmt.Errorf("events/RCONType: too many parts in event. Expected < %d got %d", len(s.Columns), len(parts))
	}

	for i, part := range parts {
		column := s.Columns[i]

		// If we're on the last part, we may have a control character we need to remove
		if isLast := len(parts)-1 == i; isLast {
			// last := part[len(part)-1:]
			part = strings.Trim(part, "\x00")
		}

		switch s.ColumnTypes[column].Kind() {
		case reflect.String:
			out[column] = strings.Trim(part, "\n")
		case reflect.Int:
			num, err := strconv.Atoi(part)
			if err != nil {
				return out, fmt.Errorf("events/RCONType: atoi(%s) on %s failed: %w", part, column, err)
			}
			out[column] = num
		case reflect.Bool:
			b, err := strconv.ParseBool(part)
			if err != nil {
				return out, fmt.Errorf("events/RCONType: strconv.ParseBool(%s) on %s failed: %w", part, column, err)
			}

			out[column] = b
		}
	}

	return out, nil
}

func ParseType(v any) RCONType {
	raw := reflect.TypeOf(v).Elem()

	// Cache to avoid a bunch of unnecessary reflection
	typeCacheMu.RLock()
	if typ, ok := typeCache[raw.Name()]; ok {
		typeCacheMu.RUnlock()
		return typ
	}
	typeCacheMu.RUnlock()

	// Need to generate, hold the RW lock
	typeCacheMu.Lock()
	defer typeCacheMu.Unlock()

	out := RCONType{
		Name:        raw.Name(),
		Columns:     make([]string, 0),
		ColumnTypes: make(map[string]reflect.Type),

		columnLookup: make(map[string]struct{}),
	}

	for i := 0; i < raw.NumField(); i++ {
		field := raw.Field(i)
		st := field.Tag.Get(expectedStructTag)
		if st == "" {
			continue // Cant do anything without a struct tag...yet
		}

		out.Columns = append(out.Columns, st)
		out.ColumnTypes[st] = field.Type
		out.columnLookup[st] = struct{}{}
	}

	typeCache[raw.Name()] = out

	return out
}

func GetAllColumns(typ any) Columns {
	rawTypeData := ParseType(typ)
	return rawTypeData.Columns
}

func UnmarshalRCON(header, data string, into any) error {
	rawTypeData := ParseType(into)

	// Extract the data based on header and input into a map
	parsedData, err := rawTypeData.Parse(header, data)
	if err != nil {
		return err
	}

	// Merge into the `into` argument
	raw := reflect.TypeOf(into).Elem()
	rawV := reflect.ValueOf(into).Elem()

	for i := 0; i < raw.NumField(); i++ {
		field := raw.Field(i)
		// Special field
		if field.Name == "LastUpdated" {
			now := time.Now()
			rawV.Field(i).Set(reflect.ValueOf(now))

			continue // Onto the next
		}

		st := field.Tag.Get(expectedStructTag)
		if st == "" {
			continue // Cant do anything without a struct tag...yet
		}

		if data, ok := parsedData[st]; ok {
			switch T := data.(type) {
			case string:
				rawV.Field(i).SetString(T)
			case bool:
				rawV.Field(i).SetBool(T)
			case int:
				rawV.Field(i).SetInt(int64(T))
			default:
				return fmt.Errorf("unexpected type of %T", data)
			}
		}
	}

	return nil
}
