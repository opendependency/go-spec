package v1

import (
	"strings"
	"testing"
)

func TestModule_Validate(t *testing.T) {
	type fields struct {
		Namespace    string
		Name         string
		Type         string
		Version      *ModuleVersion
		Annotations  map[string]string
		Dependencies []*ModuleDependency
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"is empty", fields{Namespace: "", Name: "", Type: "", Version: nil, Annotations: nil, Dependencies: nil}, true},
		{"is valid", fields{Namespace: "com.example", Name: "product", Type: "go", Version: &ModuleVersion{Name: "v1.0.0"}, Annotations: nil, Dependencies: nil}, false},
		{"has invalid namespace", fields{Namespace: "&%", Name: "product", Type: "go", Version: &ModuleVersion{Name: "v1.0.0"}, Annotations: nil, Dependencies: nil}, true},
		{"has invalid name", fields{Namespace: "com.example", Name: "&%", Type: "go", Version: &ModuleVersion{Name: "v1.0.0"}, Annotations: nil, Dependencies: nil}, true},
		{"has invalid type", fields{Namespace: "com.example", Name: "product", Type: "&%", Version: &ModuleVersion{Name: "v1.0.0"}, Annotations: nil, Dependencies: nil}, true},
		{"has invalid version", fields{Namespace: "com.example", Name: "product", Type: "go", Version: &ModuleVersion{Name: "&%"}, Annotations: nil, Dependencies: nil}, true},
		{"has invalid annotation", fields{Namespace: "com.example", Name: "product", Type: "go", Version: &ModuleVersion{Name: "v1.0.0"}, Annotations: map[string]string{"&%": ""}, Dependencies: nil}, true},
		{"has invalid dependency entry", fields{Namespace: "com.example", Name: "product", Type: "go", Version: &ModuleVersion{Name: "v1.0.0"}, Annotations: nil, Dependencies: []*ModuleDependency{{}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Module{
				Namespace:    tt.fields.Namespace,
				Name:         tt.fields.Name,
				Type:         tt.fields.Type,
				Version:      tt.fields.Version,
				Annotations:  tt.fields.Annotations,
				Dependencies: tt.fields.Dependencies,
			}
			if err := x.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleNamespace(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is empty", args{namespace: ""}, true},
		{"has min length", args{namespace: "a"}, false},
		{"has maximal length", args{namespace: strings.Repeat("a", 63)}, false},
		{"exceeds maximal length", args{namespace: strings.Repeat("a", 64)}, true},

		{"has uppercase characters", args{namespace: "A"}, true},
		{"has special characters", args{namespace: "%"}, true},

		{"starts with a number", args{namespace: "1b"}, true},
		{"starts with a letter", args{namespace: "ab"}, false},
		{"ends with a letter", args{namespace: "ab"}, false},
		{"ends with a number", args{namespace: "a0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleNamespace(tt.args.namespace); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleNamespace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is empty", args{name: ""}, true},
		{"has min length", args{name: "a"}, false},
		{"has maximal length", args{name: strings.Repeat("a", 63)}, false},
		{"exceeds maximal length", args{name: strings.Repeat("a", 64)}, true},

		{"has uppercase characters", args{name: "A"}, true},
		{"has special characters", args{name: "%"}, true},

		{"starts with a number", args{name: "1b"}, true},
		{"starts with a letter", args{name: "ab"}, false},
		{"ends with a letter", args{name: "ab"}, false},
		{"ends with a number", args{name: "a0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleType(t *testing.T) {
	type args struct {
		type_ string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is empty", args{type_: ""}, true},
		{"has min length", args{type_: "a"}, false},
		{"has maximal length", args{type_: strings.Repeat("a", 63)}, false},
		{"exceeds maximal length", args{type_: strings.Repeat("a", 64)}, true},

		{"has uppercase characters", args{type_: "A"}, true},
		{"has special characters", args{type_: "%"}, true},

		{"starts with a number", args{type_: "1b"}, true},
		{"starts with a letter", args{type_: "ab"}, false},
		{"ends with a letter", args{type_: "ab"}, false},
		{"ends with a number", args{type_: "a0"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleType(tt.args.type_); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleVersion(t *testing.T) {
	type args struct {
		moduleVersion *ModuleVersion
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is nil", args{moduleVersion: nil}, true},
		{"is invalid", args{moduleVersion: &ModuleVersion{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleVersion(tt.args.moduleVersion); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModuleVersion_Validate(t *testing.T) {
	validSchema := "my-schema"
	invalidSchema := "%&/"

	type fields struct {
		Name     string
		Schema   *string
		Replaces []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"has invalid name", fields{Name: "", Schema: nil, Replaces: nil}, true},
		{"has valid name", fields{Name: "v1.0.0", Schema: nil, Replaces: nil}, false},
		{"has invalid schema", fields{Name: "v1.0.0", Schema: &invalidSchema, Replaces: nil}, true},
		{"has valid schema", fields{Name: "v1.0.0", Schema: &validSchema, Replaces: nil}, false},
		{"has invalid replaces entry", fields{Name: "v1.0.0", Schema: nil, Replaces: []string{""}}, true},
		{"has valid replaces entry", fields{Name: "v1.1.0", Schema: nil, Replaces: []string{"v1.0.0"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &ModuleVersion{
				Name:     tt.fields.Name,
				Schema:   tt.fields.Schema,
				Replaces: tt.fields.Replaces,
			}
			if err := x.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleVersionName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is empty", args{name: ""}, true},
		{"has min length", args{name: "a"}, false},
		{"has maximal length", args{name: strings.Repeat("a", 63)}, false},
		{"exceeds maximal length", args{name: strings.Repeat("a", 64)}, true},

		{"has uppercase characters", args{name: "A"}, true},
		{"has special characters", args{name: "%"}, true},

		{"starts with a number", args{name: "1b"}, false},
		{"starts with a letter", args{name: "ab"}, false},
		{"ends with a letter", args{name: "ab"}, false},
		{"ends with a number", args{name: "a0"}, false},

		{"is valid plain version", args{name: "1.0.0"}, false},
		{"is valid plain version with identifier", args{name: "1.0.0-abc"}, false},
		{"is valid plain version with prefix v", args{name: "v1.0.0"}, false},
		{"is valid plain version with prefix v and identifier", args{name: "v1.0.0-abc"}, false},
		{"is valid date", args{name: "20210830"}, false},
		{"is valid dashed date", args{name: "2021-08-30"}, false},
		{"is valid dotted date", args{name: "2021.08.30"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleVersionName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleVersionName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleVersionSchema(t *testing.T) {
	type args struct {
		schema string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is empty", args{schema: ""}, true},
		{"has min length", args{schema: "a"}, false},
		{"has maximal length", args{schema: strings.Repeat("a", 63)}, false},
		{"exceeds maximal length", args{schema: strings.Repeat("a", 64)}, true},

		{"has uppercase characters", args{schema: "A"}, true},
		{"has special characters", args{schema: "%"}, true},

		{"starts with a number", args{schema: "1b"}, true},
		{"starts with a letter", args{schema: "ab"}, false},
		{"ends with a letter", args{schema: "ab"}, false},
		{"ends with a number", args{schema: "a0"}, false},

		{"is valid", args{schema: "my-version-schema"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleVersionSchema(tt.args.schema); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleVersionSchema() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleAnnotations(t *testing.T) {
	type args struct {
		annotations map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is nil", args{annotations: nil}, false},
		{"is empty", args{annotations: map[string]string{}}, false},

		{"has invalid key", args{annotations: map[string]string{"invalid key": ""}}, true},
		{"has invalid value", args{annotations: map[string]string{"key": strings.Repeat("a", 255)}}, true},

		{"has valid key-values", args{annotations: map[string]string{"key": "ab"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleAnnotations(tt.args.annotations); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleAnnotations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleAnnotationKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is empty", args{key: ""}, true},
		{"has min length", args{key: "a"}, false},
		{"has maximal length", args{key: strings.Repeat("a", 63)}, false},
		{"exceeds maximal length", args{key: strings.Repeat("a", 64)}, true},

		{"has uppercase characters", args{key: "A"}, true},
		{"has special characters", args{key: "%"}, true},

		{"starts with a number", args{key: "1b"}, true},
		{"starts with a letter", args{key: "ab"}, false},
		{"ends with a letter", args{key: "ab"}, false},
		{"ends with a number", args{key: "a0"}, false},

		{"is valid", args{key: "akey"}, false},
		{"is namespaces", args{key: "com.example.key"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleAnnotationKey(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleAnnotationKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleAnnotationValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"is empty", args{value: ""}, false},
		{"has maximal length", args{value: strings.Repeat("a", 253)}, false},
		{"exceeds maximal length", args{value: strings.Repeat("a", 254)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleAnnotationValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleAnnotationValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateModuleDependencies(t *testing.T) {
	type args struct {
		moduleDependencies []*ModuleDependency
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "is nil", args: args{moduleDependencies: nil}, wantErr: false},
		{name: "is empty", args: args{moduleDependencies: []*ModuleDependency{}}, wantErr: false},
		{name: "has valid entries", args: args{moduleDependencies: []*ModuleDependency{{
			Namespace: "com.example",
			Name:      "product",
			Type:      "go",
			Version:   "v1.0.0",
		}}}, wantErr: false},
		{name: "has an invalid entry", args: args{moduleDependencies: []*ModuleDependency{{
			Namespace: "com.example",
			Name:      "INVALID PRODUCT",
			Type:      "go",
			Version:   "v1.0.0",
		}}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateModuleDependencies(tt.args.moduleDependencies); (err != nil) != tt.wantErr {
				t.Errorf("validateModuleDependencies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModuleDependency_Validate(t *testing.T) {
	upstream := DependencyDirection_UPSTREAM

	type fields struct {
		Namespace string
		Name      string
		Type      string
		Version   string
		Direction *DependencyDirection
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"is empty", fields{Namespace: "", Name: "", Type: "", Version: "", Direction: nil}, true},
		{"is valid", fields{Namespace: "com.example", Name: "product", Type: "go", Version: "v1.0.0", Direction: &upstream}, false},
		{"has invalid namespace", fields{Namespace: "&%", Name: "product", Type: "go", Version: "v1.0.0", Direction: &upstream}, true},
		{"has invalid name", fields{Namespace: "com.example", Name: "PRODUCT", Type: "go", Version: "v1.0.0", Direction: &upstream}, true},
		{"has invalid type", fields{Namespace: "com.example", Name: "product", Type: "&%", Version: "v1.0.0", Direction: &upstream}, true},
		{"has invalid version", fields{Namespace: "com.example", Name: "product", Type: "go", Version: "&%", Direction: &upstream}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &ModuleDependency{
				Namespace: tt.fields.Namespace,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
				Version:   tt.fields.Version,
				Direction: tt.fields.Direction,
			}
			if err := x.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mustHaveMinMaxLength(t *testing.T) {
	type args struct {
		value  string
		minLen int
		maxLen int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"min -1, max 0: empty value", args{minLen: -1, maxLen: 0, value: ""}, true},
		{"min 0, max -1: empty value", args{minLen: 0, maxLen: -1, value: ""}, true},

		{"min 0, max 0: empty value", args{minLen: 0, maxLen: 0, value: ""}, false},
		{"min 0, max 1: empty value", args{minLen: 0, maxLen: 1, value: ""}, false},
		{"min 1, max 0: empty value", args{minLen: 1, maxLen: 0, value: ""}, true},
		{"min 1, max 1: empty value", args{minLen: 1, maxLen: 1, value: ""}, true},
		{"min 1, max 100: empty value", args{minLen: 1, maxLen: 100, value: ""}, true},

		{"min 0, max 0: one character value", args{minLen: 0, maxLen: 0, value: "a"}, true},
		{"min 0, max 1: one character value", args{minLen: 0, maxLen: 1, value: "a"}, false},
		{"min 1, max 0: one character value", args{minLen: 1, maxLen: 0, value: "a"}, true},
		{"min 1, max 1: one character value", args{minLen: 1, maxLen: 1, value: "a"}, false},

		{"min 0, max 0: two character value", args{minLen: 0, maxLen: 0, value: "ab"}, true},
		{"min 0, max 1: two character value", args{minLen: 0, maxLen: 1, value: "ab"}, true},
		{"min 0, max 2: two character value", args{minLen: 0, maxLen: 2, value: "ab"}, false},
		{"min 1, max 0: two character value", args{minLen: 1, maxLen: 0, value: "ab"}, true},
		{"min 1, max 1: two character value", args{minLen: 1, maxLen: 1, value: "ab"}, true},
		{"min 1, max 2: two character value", args{minLen: 1, maxLen: 2, value: "ab"}, false},
		{"min 1, max 100: two character value", args{minLen: 1, maxLen: 100, value: "ab"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mustHaveMinMaxLength(tt.args.value, tt.args.minLen, tt.args.maxLen); (err != nil) != tt.wantErr {
				t.Errorf("mustHaveMinMaxLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mustBeLowercaseAlphanumericDashDot(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{value: ""}, false},
		{"one character: space", args{value: " "}, true},
		{"one character: valid letter lower boundary", args{value: "a"}, false},
		{"one character: valid letter upper boundary", args{value: "z"}, false},
		{"one character: valid letter in-between boundaries", args{value: "g"}, false},
		{"one character: valid number lower boundary", args{value: "0"}, false},
		{"one character: valid number upper boundary", args{value: "9"}, false},
		{"one character: valid number in-between boundaries", args{value: "5"}, false},
		{"one character: valid dash", args{value: "-"}, false},
		{"one character: valid dot", args{value: "."}, false},
		{"one character: invalid letter character", args{value: "A"}, true},

		{"two character: valid letter lower boundary", args{value: "ab"}, false},
		{"two character: valid letter upper boundary", args{value: "yz"}, false},
		{"two character: valid letter in-between boundaries", args{value: "gh"}, false},
		{"two character: valid number lower boundary", args{value: "01"}, false},
		{"two character: valid number upper boundary", args{value: "89"}, false},
		{"two character: valid number in-between boundaries", args{value: "56"}, false},
		{"two character: valid dash", args{value: "--"}, false},
		{"two character: valid dot", args{value: ".."}, false},
		{"two character: valid mix dot and dash", args{value: ".-"}, false},
		{"two character: valid mix letter and numbers", args{value: "a0"}, false},
		{"two character: valid mix letter and dash", args{value: "a-"}, false},
		{"two character: valid mix letter and dot", args{value: "a."}, false},
		{"two character: valid mix number and dash", args{value: "0-"}, false},
		{"two character: valid mix number and dot", args{value: "0."}, false},
		{"two character: invalid letter character", args{value: "A-"}, true},

		{"three character: invalid space character", args{value: "a b"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mustBeLowercaseAlphanumericDashDot(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("mustBeLowercaseAlphanumericDashDot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mustStartWithLowercaseAlphabeticCharacter(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{value: ""}, false},
		{"one character: valid letter lower boundary", args{value: "a"}, false},
		{"one character: valid letter upper boundary", args{value: "z"}, false},
		{"one character: valid letter in-between boundaries", args{value: "g"}, false},
		{"one character: invalid number character", args{value: "0"}, true},
		{"one character: invalid character", args{value: "%"}, true},
		{"one character: invalid letter character", args{value: "A"}, true},

		{"two character: valid letter lower boundary", args{value: "a%"}, false},
		{"two character: valid letter upper boundary", args{value: "z%"}, false},
		{"two character: valid letter in-between boundaries", args{value: "g%"}, false},
		{"two character: invalid number character", args{value: "0a"}, true},
		{"two character: invalid character", args{value: "%a"}, true},
		{"two character: invalid letter character", args{value: "Aa"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mustStartWithLowercaseAlphabeticCharacter(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("mustStartWithLowercaseAlphabeticCharacter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mustStartWithLowercaseAlphanumericCharacter(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{value: ""}, false},
		{"one character: space", args{value: " "}, true},
		{"one character: valid letter lower boundary", args{value: "a"}, false},
		{"one character: valid letter upper boundary", args{value: "z"}, false},
		{"one character: valid letter in-between boundaries", args{value: "g"}, false},
		{"one character: valid number lower boundary", args{value: "0"}, false},
		{"one character: valid number upper boundary", args{value: "9"}, false},
		{"one character: valid number in-between boundaries", args{value: "5"}, false},
		{"one character: invalid character", args{value: "%"}, true},
		{"one character: invalid letter character", args{value: "A"}, true},

		{"two character: valid letter lower boundary", args{value: "a%"}, false},
		{"two character: valid letter upper boundary", args{value: "z%"}, false},
		{"two character: valid letter in-between boundaries", args{value: "g%"}, false},
		{"two character: valid number lower boundary", args{value: "0%"}, false},
		{"two character: valid number upper boundary", args{value: "9%"}, false},
		{"two character: valid number in-between boundaries", args{value: "5%"}, false},
		{"two character: invalid character", args{value: "%a"}, true},
		{"two character: invalid letter character", args{value: "Aa"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mustStartWithLowercaseAlphanumericCharacter(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("mustStartWithLowercaseAlphanumericCharacter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mustEndWithLowercaseAlphanumericCharacter(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{value: ""}, false},
		{"one character: space", args{value: " "}, true},
		{"one character: valid letter lower boundary", args{value: "a"}, false},
		{"one character: valid letter upper boundary", args{value: "z"}, false},
		{"one character: valid letter in-between boundaries", args{value: "g"}, false},
		{"one character: valid number lower boundary", args{value: "0"}, false},
		{"one character: valid number upper boundary", args{value: "9"}, false},
		{"one character: valid number in-between boundaries", args{value: "5"}, false},
		{"one character: invalid character", args{value: "%"}, true},
		{"one character: invalid letter character", args{value: "A"}, true},

		{"two character: valid letter lower boundary", args{value: "%a"}, false},
		{"two character: valid letter upper boundary", args{value: "%z"}, false},
		{"two character: valid letter in-between boundaries", args{value: "%g"}, false},
		{"two character: valid number lower boundary", args{value: "%0"}, false},
		{"two character: valid number upper boundary", args{value: "%9"}, false},
		{"two character: valid number in-between boundaries", args{value: "%5"}, false},
		{"two character: invalid character", args{value: "a%"}, true},
		{"two character: invalid letter character", args{value: "aA"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := mustEndWithLowercaseAlphanumericCharacter(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("mustEndWithLowercaseAlphanumericCharacter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
