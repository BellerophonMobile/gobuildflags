package goflagbuilder

import (
	"bytes"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mystruct struct {
	FieldA string `help:"Field A"`
	FieldB int
}

type myotherstruct struct {
	Grid     uint64
	Fraction float64

	Attrs map[string]string
}

type mystruct2 struct {
	Name  string
	Index int

	DoStuff bool

	Location myotherstruct
}

type mystruct3 struct {
	Name  string
	Index int

	Location *myotherstruct
}

type mystruct4 struct {
	Temp  int64
	Check uint
	Files []string
}

type mystruct5 struct {
	Getter testGetter
}

type configStruct struct {
	A struct {
		Name  string `help:"Name"`
		Value string `help:"Value"`
	}
	B struct {
		Foo int64
		Bar uint
	}
}

type testGetter struct {
	value string
}

func (t *testGetter) Set(s string) error {
	t.value = s
	return nil
}

func (t *testGetter) Get() interface{} { return t.value }
func (t *testGetter) String() string   { return t.value }

func TestInto_Invalid(t *testing.T) {
	suite := []struct {
		name string
		conf interface{}
		err  string
	}{
		{"nil", nil, "cannot build flags from nil"},
		{"string", "Banana", "cannot build flags from type string for prefix ''"},
		{"int", 7, "cannot build flags from type int for prefix ''"},
		{"float", 7.0, "cannot build flags from type float64 for prefix ''"},
		{
			name: "struct",
			conf: mystruct{"Banana", 7},
			err:  "value of type string at FieldA cannot be set",
		},
		{
			name: "map to struct",
			conf: map[string]interface{}{"MyStruct": mystruct{"Banana", 7}},
			err:  "value of type string at MyStruct.FieldA cannot be set",
		},
		{
			name: "map without string keys",
			conf: map[int]interface{}{10: mystruct{"Banana", 7}},
			err:  "map key must be string, got int for prefix ''",
		},
	}

	for _, item := range suite {
		t.Run(item.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet(item.name, flag.ContinueOnError)
			err := Into(flagSet, item.conf)

			if err == nil {
				t.Error("did not get expected error:", item.err)
			} else {
				actualError := err.Error()
				if item.err != actualError {
					t.Errorf("expected error did not match actual.\nExpected: %s\n  Actual: %s", item.err, actualError)
				}
			}
		})
	}
}

type expectedVariable struct {
	value interface{}
	usage string
}

func TestInto(t *testing.T) {
	suite := []struct {
		name string
		conf interface{}
		args []string
		vars map[string]expectedVariable
		help string
	}{
		{
			name: "empty map",
			conf: map[string]int{},
			args: []string{},
			vars: map[string]expectedVariable{},
		},
		{
			name: "map to int",
			conf: map[string]int{"Banana": 7},
			args: []string{"-Banana", "10"},
			vars: map[string]expectedVariable{
				"Banana": {value: 10},
			},
			help: "  -Banana\n    \t (default 7)\n",
		},
		{
			name: "map to struct ptr",
			conf: map[string]interface{}{"MyStruct": &mystruct{}},
			args: []string{"-MyStruct.FieldA", "asdf", "-MyStruct.FieldB", "12"},
			vars: map[string]expectedVariable{
				"MyStruct.FieldA": {value: "asdf", usage: "Field A"},
				"MyStruct.FieldB": {value: 12},
			},
			help: "  -MyStruct.FieldA\n    \tField A\n  -MyStruct.FieldB\n    \t (default 0)\n",
		},
		{
			name: "struct ptr",
			conf: &mystruct{},
			args: []string{"-FieldA", "foo", "-FieldB", "21"},
			vars: map[string]expectedVariable{
				"FieldA": {value: "foo", usage: "Field A"},
				"FieldB": {value: 21},
			},
			help: "  -FieldA\n    \tField A\n  -FieldB\n    \t (default 0)\n",
		},
		{
			name: "nested struct",
			conf: &mystruct2{},
			args: []string{"-Name", "foo", "-Index", "10", "-DoStuff", "-Location.Grid", "2048", "-Location.Fraction", "3.14"},
			vars: map[string]expectedVariable{
				"Name":              {value: "foo"},
				"Index":             {value: 10},
				"DoStuff":           {value: true},
				"Location.Grid":     {value: uint64(2048)},
				"Location.Fraction": {value: 3.14},
			},
			help: "  -DoStuff\n    \t (default false)\n  -Index\n    \t (default 0)\n  -Location.Fraction\n    \t (default 0)\n  -Location.Grid\n    \t (default 0)\n  -Name\n    \t\n",
		},
		{
			name: "nested struct ptr",
			conf: &mystruct3{Location: &myotherstruct{}},
			args: []string{"-Name", "bar", "-Index", "20", "-Location.Grid", "1000", "-Location.Fraction", "2.71"},
			vars: map[string]expectedVariable{
				"Name":              {value: "bar"},
				"Index":             {value: 20},
				"Location.Grid":     {value: uint64(1000)},
				"Location.Fraction": {value: 2.71},
			},
			help: "  -Index\n    \t (default 0)\n  -Location.Fraction\n    \t (default 0)\n  -Location.Grid\n    \t (default 0)\n  -Name\n    \t\n",
		},
		{
			name: "struct with nested map",
			conf: &myotherstruct{Attrs: map[string]string{"Foo": "Bar"}},
			args: []string{"-Grid", "12", "-Fraction", "1.23", "-Attrs.Foo", "AAA"},
			vars: map[string]expectedVariable{
				"Grid":      {value: uint64(12)},
				"Fraction":  {value: 1.23},
				"Attrs.Foo": {value: "AAA"},
			},
			help: "  -Attrs.Foo\n    \t (default Bar)\n  -Fraction\n    \t (default 0)\n  -Grid\n    \t (default 0)\n",
		},
		{
			name: "struct with nil pointer",
			conf: &mystruct3{},
			args: []string{"-Index", "10", "-Location.Fraction", "3.14"},
			vars: map[string]expectedVariable{
				"Name":              {value: ""},
				"Index":             {value: 10},
				"Location.Grid":     {value: uint64(0)},
				"Location.Fraction": {value: 3.14},
			},
			help: "  -Index\n    \t (default 0)\n  -Location.Fraction\n    \t (default 0)\n  -Location.Grid\n    \t (default 0)\n  -Name\n    \t\n",
		},
		{
			name: "struct with slice",
			conf: &mystruct4{},
			args: []string{"-Temp", "-10", "-Check", "5", "-Files", "foo.log", "-Files", "bar.txt"},
			vars: map[string]expectedVariable{
				"Temp":  {value: int64(-10)},
				"Check": {value: uint(5)},
				"Files": {value: []string{"foo.log", "bar.txt"}},
			},
			help: "  -Check\n    \t (default 0)\n  -Files\n    \t (default [])\n  -Temp\n    \t (default 0)\n",
		},
		{
			name: "struct with getter",
			conf: &mystruct5{},
			args: []string{"-Getter", "Foo"},
			vars: map[string]expectedVariable{
				"Getter": {value: "Foo"},
			},
			help: "  -Getter\n    \t\n",
		},
		{
			name: "struct with nested anoymous structs",
			conf: &configStruct{},
			args: []string{"-A.Name", "asdf"},
			vars: map[string]expectedVariable{
				"A.Name":  {value: "asdf", usage: "Name"},
				"A.Value": {value: "", usage: "Value"},
				"B.Foo":   {value: int64(0)},
				"B.Bar":   {value: uint(0)},
			},
			help: "  -A.Name\n    \tName\n  -A.Value\n    \tValue\n  -B.Bar\n    \t (default 0)\n  -B.Foo\n    \t (default 0)\n",
		},
	}

	for _, item := range suite {
		t.Run(item.name, func(t *testing.T) {
			flagSet := flag.NewFlagSet(item.name, flag.ContinueOnError)

			if err := Into(flagSet, item.conf); err != nil {
				t.Error("unexpected error:", err)
				return
			}

			if err := flagSet.Parse(item.args); err != nil {
				t.Error("error parsing args:", err)
				return
			}

			flagSet.VisitAll(func(f *flag.Flag) {
				if _, ok := item.vars[f.Name]; ok {
					return
				}
				t.Error("unexpected variable:", f.Name)
			})

			for name, expected := range item.vars {
				f := flagSet.Lookup(name)
				if f == nil {
					t.Error("expected variable was not found:", name)
					return
				}

				if expected.usage != f.Usage {
					t.Errorf("usage doesn't match\nexpected: %s\n  actual: %s", expected.usage, f.Usage)
				}

				getter, ok := f.Value.(flag.Getter)
				if !ok {
					t.Fatal("value not getter?")
					return
				}

				assert.Equal(t, expected.value, getter.Get(), "values for %s not equal", name)
			}

			var buf bytes.Buffer

			flagSet.SetOutput(&buf)
			flagSet.PrintDefaults()

			assert.Equal(t, item.help, buf.String())
		})
	}
}
