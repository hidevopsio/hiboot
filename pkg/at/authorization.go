package at


// RequiresAuthentication is the annotation that annotate the method for authorization
type RequiresAuthentication struct {
	Annotation

	BaseAnnotation
}


// RequiresRoles is the annotation that annotate the method for requires roles
type RequiresRoles struct {
	Annotation

	BaseAnnotation
}


// RequiresPermissions  is the annotation that annotate the method for requires permissions
type RequiresPermissions  struct {
	Annotation

	BaseAnnotation
}

// RequiresUser  is the annotation that annotate the method for requires users
type RequiresUser  struct {
	Annotation

	BaseAnnotation
}
