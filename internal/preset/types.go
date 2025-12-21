package preset

import "time"

type Package struct {
	Version      string            `yaml:"version" json:"version"`
	Name         string            `yaml:"name" json:"name"`
	Region       string            `yaml:"region" json:"region"`
	Description  string            `yaml:"description,omitempty" json:"description,omitempty"`
	Tags         []string          `yaml:"tags,omitempty" json:"tags,omitempty"`
	Dependencies map[string]string `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`

	EntitiesFiles         []EntityFile `yaml:"-" json:"entities_files"`
	EntitiesCount         int          `yaml:"-" json:"entities_count"`
	EntitiesTotalSize     int64        `yaml:"-" json:"entities_total_size"`
	EntitiesStructureHash uint32       `yaml:"-" json:"entities_structure_hash"`

	Entities []RowEntity `yaml:"-" json:"-"`
}

type EntityFile struct {
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	ContentHash uint32    `json:"content_hash,omitempty"`
}

type ProcessedEntity struct {
	File        EntityFile     // Метаданные файла
	ContentHash uint32         // Хеш содержимого
	JSONData    []byte         // YAML → JSON (готовый для сохранения)
	ParsedData  map[string]any // ТОЛЬКО для быстрой валидации
	Schema      map[string]any `json:"schema"` // JSON Schema
	Errors      []string       // Ошибки валидации
	FatalError  error          // Фатальная ошибка чтения/конвертации
}

type EntityOutput struct {
	Metadata   EntityMetadata   `json:"metadata"`
	ParsedData map[string]any   `json:"parsed_data"`
	JSONData   string           `json:"json_data,omitempty"`
	Schema     map[string]any   `json:"schema"`
	Validation ValidationResult `json:"validation"`
}

type ValidationResult struct {
	IsValid    bool     `json:"is_valid"`
	HasFatal   bool     `json:"has_fatal"`
	ErrorCount int      `json:"error_count"`
	Errors     []string `json:"errors,omitempty"`
	FatalError string   `json:"fatal_error,omitempty"`
}

type EntityMetadata struct {
	Module      string    `json:"module"`
	Object      string    `json:"object"`
	Property    string    `json:"property"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	SourceFile  string    `json:"source_file"`
	FileSize    int64     `json:"file_size"`
	ModTime     time.Time `json:"mod_time"`
	ContentHash string    `json:"content_hash"`
	ProcessedAt time.Time `json:"processed_at"`
}

type RowEntity struct {
	Module   string     `json:"module"`
	Object   string     `json:"object"`
	Property string     `json:"property"`
	Code     string     `json:"code"`
	Name     string     `json:"name"`
	Fields   []RowField `json:"fields"`
}

type RowField struct {
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Type     FieldType `json:"type"`
	Pattern  *string   `json:"pattern,omitempty"`
	Required bool      `json:"required"`
	Min      *int      `json:"min,omitempty"`
	Max      *int      `json:"max,omitempty"`

	DefaultString  *string  `json:"default_string,omitempty"`
	DefaultNumber  *float64 `json:"default_number,omitempty"`
	DefaultInteger *int     `json:"default_integer,omitempty"`
	DefaultBoolean *bool    `json:"default_boolean,omitempty"`
	DefaultEnum    *string  `json:"default_enum,omitempty"`

	EnumValues []string `json:"enum_values,omitempty"`

	MultipleOf *float64 `json:"multiple_of,omitempty"`
}

type FieldType string

const (
	TypeString  FieldType = "string"
	TypeNumber  FieldType = "number"
	TypeInteger FieldType = "integer"
	TypeBoolean FieldType = "boolean"
	TypeEnum    FieldType = "enum"
)
