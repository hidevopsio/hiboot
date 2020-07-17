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

	AtLogical Logical `json:"-" at:"logical"`
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
}

// RequiresUser  is the annotation that annotate the method for requires users
type RequiresUser  struct {
	Annotation

	BaseAnnotation
}

// RequiresData  is the annotation that annotate the method for requires data access
type RequiresData  struct {
	Annotation

	RequiresLogical

	AtPermission string `json:"-" at:"permission"`
	AtAction string `json:"-" at:"action"`
	AtCondition string `json:"-" at:"condition"`
}

