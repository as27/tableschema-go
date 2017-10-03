package schema

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestDefaultValues(t *testing.T) {
	data := []struct {
		Desc  string
		JSON  string
		Field Field
	}{
		{
			"Default Values",
			`{"name":"n1"}`,
			Field{Name: "n1", Type: defaultFieldType, Format: defaultFieldFormat, TrueValues: defaultTrueValues, FalseValues: defaultFalseValues,
				DecimalChar: defaultDecimalChar, GroupChar: defaultGroupChar, BareNumber: defaultBareNumber},
		},
		{
			"Overrinding default values",
			`{"name":"n2","type":"t2","format":"f2","falseValues":["f2"],"trueValues":["t2"]}`,
			Field{Name: "n2", Type: "t2", Format: "f2", TrueValues: []string{"t2"}, FalseValues: []string{"f2"},
				DecimalChar: defaultDecimalChar, GroupChar: defaultGroupChar, BareNumber: defaultBareNumber},
		},
	}
	for _, d := range data {
		var f Field
		if err := json.Unmarshal([]byte(d.JSON), &f); err != nil {
			t.Errorf("err want:nil got:%q", err)
		}
		if !reflect.DeepEqual(f, d.Field) {
			t.Errorf("[%s] want:%+v got:%+v", d.Desc, d.Field, f)
		}
	}
}

func TestField_Decode(t *testing.T) {
	data := []struct {
		Desc     string
		Value    string
		Field    Field
		Expected interface{}
	}{
		{"Integer", "42", Field{Type: IntegerType}, int64(42)},
		{"String_URI", "http:/frictionlessdata.io", Field{Type: StringType, Format: "uri"}, "http:/frictionlessdata.io"},
		{"Boolean_TrueValues", "1", Field{Type: BooleanType, TrueValues: []string{"1"}}, true},
		{"Boolean_FalseValues", "0", Field{Type: BooleanType, FalseValues: []string{"0"}}, false},
		{"Number", "42.5", Field{Type: NumberType}, 42.5},
		{"Date_NoFormat", "2015-10-15", Field{Type: DateType}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"Date_DefaultFormat", "2015-10-15", Field{Type: DateType, Format: defaultFieldFormat}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"Date_CustomFormat", "15/10/2015", Field{Type: DateType, Format: "%d/%m/%Y"}, time.Date(2015, time.October, 15, 0, 0, 0, 0, time.UTC)},
		{"Time_NoFormat", "10:10:10", Field{Type: TimeType}, time.Date(0000, time.January, 01, 10, 10, 10, 00, time.UTC)},
		{"Time_DefaultFormat", "10:10:10", Field{Type: TimeType, Format: defaultFieldFormat}, time.Date(0000, time.January, 01, 10, 10, 10, 00, time.UTC)},
		{"Time_CustomFormat", "10-10-10", Field{Type: TimeType, Format: "%H-%M-%S"}, time.Date(0000, time.January, 01, 10, 10, 10, 00, time.UTC)},
		{"YearMonth", "2017-08", Field{Type: YearMonthType}, time.Date(2017, time.August, 01, 00, 00, 00, 00, time.UTC)},
		{"Year", "2017", Field{Type: YearType}, time.Date(2017, time.January, 01, 00, 00, 00, 00, time.UTC)},
		{"DateTime_NoFormat", "2008-09-15T10:53:00Z", Field{Type: DateTimeType}, time.Date(2008, time.September, 15, 10, 53, 00, 00, time.UTC)},
		{"DateTime_DefaultFormat", "2008-09-15T10:53:00Z", Field{Type: DateTimeType, Format: defaultFieldFormat}, time.Date(2008, time.September, 15, 10, 53, 00, 00, time.UTC)},
		{"Duration", "P2H", Field{Type: DurationType}, 2 * time.Hour},
		{"GeoPoint", "90,45", Field{Type: GeoPointType}, GeoPoint{90, 45}},
		{"Any", "10", Field{Type: AnyType}, "10"},
	}
	for _, d := range data {
		t.Run(d.Desc, func(t *testing.T) {
			c, err := d.Field.Decode(d.Value)
			if err != nil {
				t.Fatalf("err want:nil got:%s", err)
			}
			if c != d.Expected {
				t.Errorf("val want:%v, got:%v", d.Expected, c)
			}
		})
	}
	t.Run("Object_Success", func(t *testing.T) {
		f := Field{Type: ObjectType}
		obj, err := f.Decode(`{"name":"foo"}`)
		if err != nil {
			t.Fatalf("err want:nil got:%s", err)
		}
		objMap, ok := obj.(map[string]interface{})
		if !ok {
			t.Errorf("want:true got:false")
		}
		if len(objMap) != 1 {
			t.Errorf("want:1 got:%d", len(objMap))
		}
		if objMap["name"] != "foo" {
			t.Errorf("val want:map[name:foo], got:%v", objMap)
		}
	})
	t.Run("Object_Failure", func(t *testing.T) {
		f := Field{Type: ObjectType}
		_, err := f.Decode(`{"name"}`)
		if err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("Array_Success", func(t *testing.T) {
		f := Field{Type: ArrayType}
		obj, err := f.Decode(`["foo"]`)
		if err != nil {
			t.Fatalf("err want:nil got:%s", err)
		}
		arr, ok := obj.([]interface{})
		if !ok {
			t.Errorf("want:true got:false")
		}
		if len(arr) != 1 {
			t.Errorf("want:1 got:%d", len(arr))
		}
		if arr[0] != "foo" {
			t.Errorf("val want:foo, got:%v", arr)
		}
	})
	t.Run("Array_Failure", func(t *testing.T) {
		f := Field{Type: ArrayType}
		_, err := f.Decode(`{"name":"foo"}`)
		if err == nil {
			t.Fatalf("err want:err got:nil")
		}
	})
	t.Run("InvalidDate", func(t *testing.T) {
		data := []struct {
			desc  string
			field Field
			value string
		}{
			{"InvalidFormat_Any", Field{Type: DateType, Format: "any"}, "2015-10-15"},
			{"InvalidFormat_Strftime", Field{Type: DateType, Format: "Fooo"}, "2015-10-15"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				if _, err := d.field.Decode(d.value); err == nil {
					t.Errorf("want:err got:nil")
				}
			})
		}
	})
	t.Run("InvalidFieldType", func(t *testing.T) {
		f := Field{Type: "invalidType"}
		if _, err := f.Decode("42"); err == nil {
			t.Errorf("err want:err got:nil")
		}
	})
	t.Run("Constraints", func(t *testing.T) {
		t.Run("Required", func(t *testing.T) {
			f := Field{Type: StringType, Constraints: Constraints{Required: true}, MissingValues: map[string]struct{}{"NA": struct{}{}}}
			if _, err := f.Decode("NA"); err == nil {
				t.Fatalf("err want:err got:nil")
			}
		})
	})
}

func TestUnmarshalJSON_InvalidField(t *testing.T) {
	var f Field
	if err := json.Unmarshal([]byte("{Foo:1}"), &f); err == nil {
		t.Errorf("want:err got:nil")
	}
}

func TestTestString(t *testing.T) {
	f := Field{Type: "integer"}
	if !f.TestString("42") {
		t.Errorf("want:true, got:false")
	}
	if f.TestString("boo") {
		t.Errorf("want:false, got:true")
	}
}

func TestField_Encode(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		data := []struct {
			desc  string
			field Field
			value interface{}
			want  string
		}{
			{"Int", Field{Type: IntegerType}, 1, "1"},
			{"Number", Field{Type: NumberType}, 1.0, "1"},
			{"IntNumberImplicitCast", Field{Type: NumberType}, 100, "100"},
			{"NumberToIntImplicitCast", Field{Type: IntegerType}, 100.5, "100"},
			{"Boolean", Field{Type: BooleanType}, true, "true"},
			{"Duration", Field{Type: DurationType}, 1 * time.Second, "P0Y0M0DT1S"},
			{"GeoPoint", Field{Type: GeoPointType}, "10,10", "10,10"},
			{"String", Field{Type: StringType}, "foo", "foo"},
			{"Array", Field{Type: ArrayType}, []string{"foo"}, "[foo]"},
			{"Date", Field{Type: DateType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"Year", Field{Type: YearType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"YearMonth", Field{Type: YearMonthType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"DateTime", Field{Type: DateTimeType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"Date", Field{Type: DateType}, time.Unix(1, 0), "1970-01-01T00:00:01Z"},
			{"Object", Field{Type: ObjectType}, eoStruct{Name: "Foo"}, `{"name":"Foo"}`},
			{"Any", Field{Type: AnyType}, "10", "10"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				got, err := d.field.Encode(d.value)
				if err != nil {
					t.Fatalf("err want:nil got:%q", err)
				}
				if d.want != got {
					t.Fatalf("val want:%s got:%s", d.want, got)
				}
			})
		}
	})
	t.Run("Error", func(t *testing.T) {
		data := []struct {
			desc  string
			field Field
			value interface{}
		}{
			{"StringToIntCast", Field{Type: IntegerType}, "1.5"},
			{"StringToNumberCast", Field{Type: NumberType}, "1.5"},
			{"InvalidType", Field{Type: "Boo"}, "1"},
		}
		for _, d := range data {
			t.Run(d.desc, func(t *testing.T) {
				if _, err := d.field.Encode(d.value); err == nil {
					t.Errorf("want:err got:nil")
				}
			})
		}
	})
}
