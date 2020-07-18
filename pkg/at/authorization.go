package at

type Logical string

const (
	AND = "and"
	OR = "or"
)

// RequiresLogical  is the annotation that annotate the method for requires logical
type RequiresLogical  struct {
	Annotation

	BaseAnnotation

	AtLogical Logical `json:"-" at:"logical" value:"and"` // default value is and
}

// RequiresAuthentication is the annotation that annotate the method for authorization
type RequiresAuthentication struct {
	Annotation

	BaseAnnotation
}


// RequiresRoles is the annotation that annotate the method for requires roles
type RequiresRoles struct {
	Annotation

	RequiresLogical
}


// RequiresPermissions  is the annotation that annotate the method for requires permissions
type RequiresPermissions  struct {
	Annotation

	RequiresLogical

	AtValues []string `at:"values" json:"-"`	// `values:"user:read,team:read"`

	AtType string `json:"-" at:"type"`	// `type:"pagination"`

	AtIn []string `json:"-" at:"in" `	// `in:"page,per_page"`

	AtOut []string `json:"-" at:"out"` 	// `out:"expr"` <where in (1,2,3)>
}

// RequiresUser  is the annotation that annotate the method for requires users
type RequiresUser  struct {
	Annotation

	BaseAnnotation
}


// RequiresGuest  is the annotation that annotate the method for requires guest
type RequiresGuest  struct {
	Annotation

	BaseAnnotation
}

