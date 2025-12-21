package preset

type Package struct {
	Version      string            `yaml:"version" json:"version"`
	Name         string            `yaml:"name" json:"name"`
	Region       string            `yaml:"region" json:"region"`
	Description  string            `yaml:"description,omitempty" json:"description,omitempty"`
	Tags         []string          `yaml:"tags,omitempty" json:"tags,omitempty"`
	Dependencies map[string]string `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
	Entities     []Entity          `yaml:"-" json:"entities"`
}

type Config struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Region  string `json:"string"`
}

type Entity struct {
	Module   string  `json:"module"`
	Object   string  `json:"object"`
	Property string  `json:"property"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Fields   []Field `json:"fields"`
}

type Field struct {
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
