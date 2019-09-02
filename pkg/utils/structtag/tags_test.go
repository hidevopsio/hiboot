package structtag

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func TestParse(t *testing.T) {
	test := []struct {
		name    string
		tag     string
		exp     []*Tag
		invalid bool
	}{
		{
			name:    "wrong tag",
			tag:     `v:'1`,
			invalid: true,
		},
		{
			name:    "wrong tag",
			tag:     `v:"1`,
			invalid: true,
		},
		{
			name: "empty tag",
			tag:  "",
		},
		{
			name:    "wrong tag",
			tag:     "\"",
			invalid: true,
		},
		{
			name:    "empty tag",
			tag:     `value:a"`,
			invalid: true,
		},
		{
			name:    "tag with one key (invalid)",
			tag:     "json",
			invalid: true,
		},
		{
			name: "tag with one key (valid)",
			tag:  `json:""`,
			exp: []*Tag{
				{
					Key: "json",
				},
			},
		},
		{
			name: "tag with one key and dash name",
			tag:  `json:"-"`,
			exp: []*Tag{
				{
					Key:  "json",
					Name: "-",
				},
			},
		},
		{
			name: "tag with key and name",
			tag:  `json:"foo"`,
			exp: []*Tag{
				{
					Key:  "json",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with key and name",
			tag:  `json:"foo" `,
			exp: []*Tag{
				{
					Key:  "json",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with key and name",
			tag:  `json:"\\foo" `,
			exp: []*Tag{
				{
					Key:  "json",
					Name: `\\foo`,
				},
			},
		},
		{
			name: "tag with key, name and option",
			tag:  `json:"foo,omitempty"`,
			exp: []*Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{"omitempty"},
				},
			},
		},
		{
			name: "tag with multiple keys",
			tag:  `json:"" hcl:""`,
			exp: []*Tag{
				{
					Key: "json",
				},
				{
					Key: "hcl",
				},
			},
		},
		{
			name: "tag with multiple keys and names",
			tag:  `json:"foo" hcl:"foo"`,
			exp: []*Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with multiple keys and names",
			tag:  `json:"foo" hcl:"foo"`,
			exp: []*Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with multiple keys and different names",
			tag:  `json:"foo" hcl:"bar"`,
			exp: []*Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "bar",
				},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo,omitempty" structs:"bar,omitnested"`,
			exp: []*Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{"omitempty"},
				},
				{
					Key:     "structs",
					Name:    "bar",
					Options: []string{"omitnested"},
				},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo" structs:"bar,omitnested" hcl:"-"`,
			exp: []*Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:     "structs",
					Name:    "bar",
					Options: []string{"omitnested"},
				},
				{
					Key:  "hcl",
					Name: "-",
				},
			},
		},
		{
			name: "tag with quoted name",
			tag:  `json:"foo,bar:\"baz\""`,
			exp: []*Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{`bar:\"baz\"`},
				},
			},
		},
	}

	for _, ts := range test {
		t.Run(ts.name, func(t *testing.T) {
			tags, err := Parse(ts.tag)
			invalid := err != nil

			if invalid != ts.invalid {
				t.Errorf("invalid case\n\twant: %+v\n\tgot : %+v\n\terr : %s", ts.invalid, invalid, err)
			}

			if invalid {
				return
			}

			got := tags.Tags()

			if !reflect.DeepEqual(ts.exp, got) {
				t.Errorf("parse\n\twant: %#v\n\tgot : %#v", ts.exp, got)
			}
		})
	}
}

func TestGoString(t *testing.T) {

	t.Run("should get all go string", func(t *testing.T) {
		tag := &Tag{
			Key:     "json",
			Name:    "bar",
			Options: []string{},
		}
		gs := tag.GoString()
		assert.Equal(t, `{
		Key:    'json',
		Name:   'bar',
		Option: '',
	}`, gs)
	})

	t.Run("should get all go string with nil option", func(t *testing.T) {
		tag := &Tag{
			Key:     "json",
			Name:    "bar",
			Options: nil,
		}
		gs := tag.GoString()
		assert.Equal(t, `{
		Key:    'json',
		Name:   'bar',
		Option: 'nil',
	}`, gs)
	})
}

func TestTags_Get(t *testing.T) {
	tag := `json:"foo,omitempty" structs:"bar,omitnested"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("go string", func(t *testing.T) {
		keys := tags.Keys()
		assert.Equal(t, "json", keys[0])
	})

	found, err := tags.Get("json")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("not exist", func(t *testing.T) {
		_, err := tags.Get("xxx")
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("String", func(t *testing.T) {
		want := `json:"foo,omitempty"`
		if found.String() != want {
			t.Errorf("get\n\twant: %#v\n\tgot : %#v", want, found.String())
		}
	})
	t.Run("Value", func(t *testing.T) {
		want := `foo,omitempty`
		if found.Value() != want {
			t.Errorf("get\n\twant: %#v\n\tgot : %#v", want, found.Value())
		}
	})
}

func TestTags_Set(t *testing.T) {
	tag := `json:"foo,omitempty" structs:"bar,omitnested"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	err = tags.Set(&Tag{
		Key:     "json",
		Name:    "bar",
		Options: []string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	found, err := tags.Get("json")
	if err != nil {
		t.Fatal(err)
	}

	want := `json:"bar"`
	if found.String() != want {
		t.Errorf("set\n\twant: %#v\n\tgot : %#v", want, found.String())
	}
}

func TestTags_Set_Append(t *testing.T) {
	tag := `json:"foo,omitempty"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	err = tags.Set(&Tag{
		Key:     "structs",
		Name:    "bar",
		Options: []string{"omitnested"},
	})
	if err != nil {
		t.Fatal(err)
	}

	found, err := tags.Get("structs")
	if err != nil {
		t.Fatal(err)
	}

	want := `structs:"bar,omitnested"`
	if found.String() != want {
		t.Errorf("set append\n\twant: %#v\n\tgot : %#v", want, found.String())
	}

	wantFull := `json:"foo,omitempty" structs:"bar,omitnested"`
	if tags.String() != wantFull {
		t.Errorf("set append\n\twant: %#v\n\tgot : %#v", wantFull, tags.String())
	}
}

func TestTags_Set_KeyDoesNotExist(t *testing.T) {
	tag := `json:"foo,omitempty" structs:"bar,omitnested"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	err = tags.Set(&Tag{
		Key:     "",
		Name:    "bar",
		Options: []string{},
	})
	if err == nil {
		t.Fatal("setting tag with a nonexisting key should error")
	}

	if err != errKeyNotSet {
		t.Errorf("set\n\twant: %#v\n\tgot : %#v", errTagKeyMismatch, err)
	}
}

func TestTags_Delete(t *testing.T) {
	tag := `json:"foo,omitempty" structs:"bar,omitnested" hcl:"-"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	tags.Delete("structs")
	if tags.Len() != 2 {
		t.Fatalf("tag length should be 2, have %d", tags.Len())
	}

	found, err := tags.Get("json")
	if err != nil {
		t.Fatal(err)
	}

	want := `json:"foo,omitempty"`
	if found.String() != want {
		t.Errorf("delete\n\twant: %#v\n\tgot : %#v", want, found.String())
	}

	wantFull := `json:"foo,omitempty" hcl:"-"`
	if tags.String() != wantFull {
		t.Errorf("delete\n\twant: %#v\n\tgot : %#v", wantFull, tags.String())
	}
}

func TestTags_DeleteOptions(t *testing.T) {
	tag := `json:"foo,omitempty" structs:"bar,omitnested,omitempty" hcl:"-"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	tags.DeleteOptions("json", "omitempty")

	want := `json:"foo" structs:"bar,omitnested,omitempty" hcl:"-"`
	if tags.String() != want {
		t.Errorf("delete option\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}

	tags.DeleteOptions("structs", "omitnested")
	want = `json:"foo" structs:"bar,omitempty" hcl:"-"`
	if tags.String() != want {
		t.Errorf("delete option\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}
}

func TestTags_AddOption(t *testing.T) {
	tag := `json:"foo" structs:"bar,omitempty" hcl:"-"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	tags.AddOptions("json", "omitempty")

	want := `json:"foo,omitempty" structs:"bar,omitempty" hcl:"-"`
	if tags.String() != want {
		t.Errorf("add options\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}

	// this shouldn't change anything
	tags.AddOptions("structs", "omitempty")

	want = `json:"foo,omitempty" structs:"bar,omitempty" hcl:"-"`
	if tags.String() != want {
		t.Errorf("add options\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}

	// this should append to the existing
	tags.AddOptions("structs", "omitnested", "flatten")
	want = `json:"foo,omitempty" structs:"bar,omitempty,omitnested,flatten" hcl:"-"`
	if tags.String() != want {
		t.Errorf("add options\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}
}

func TestTags_String(t *testing.T) {
	tag := `json:"foo" structs:"bar,omitnested" hcl:"-"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	if tags.String() != tag {
		t.Errorf("string\n\twant: %#v\n\tgot : %#v", tag, tags.String())
	}
}

func TestTags_Sort(t *testing.T) {
	tag := `json:"foo" structs:"bar,omitnested" hcl:"-"`

	tags, err := Parse(tag)
	if err != nil {
		t.Fatal(err)
	}

	sort.Sort(tags)

	want := `hcl:"-" json:"foo" structs:"bar,omitnested"`
	if tags.String() != want {
		t.Errorf("string\n\twant: %#v\n\tgot : %#v", want, tags.String())
	}
}

func TestTag_String(t *testing.T) {
	tags := new(Tags)
	s := tags.String()
	assert.Equal(t, "", s)
}

type ty struct {
	Foo string `tag:"foo,bar:\"baz\""`
}

//func TestParse2(t *testing.T) {
//	tagString := string(reflect.TypeOf(ty{}).Field(0).Tag)
//	tags, err := Parse(tagString)
//	require.NoError(t, err)
//	tag, err := tags.Get("tag")
//	require.NoError(t, err)
//	require.Equal(t, `bar:"baz"`, tag.Options[0])
//	require.Equal(t, tagString, tags.String())
//}