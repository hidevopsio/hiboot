package at

// Scope annotation in hiboot is used to define a prototype scope for an instance. In Hiboot,
// an instance can have different scopes which determine the lifecycle and visibility of that
// instance. The prototype scope means that a new instance of the object will be created every
// time it is requested from the Hiboot container.
//
//	type Example struct {
//	  at.Scope `value:"singleton"` // singleton
//	  ...
//	}
//
//	type Example struct {
//	  at.Scope `value:"prototype"` // prototype
//	  ...
//	}

type Scope struct {
	Annotation `json:"-"`

	BaseAnnotation
}

// ContextAware is the annotation that has the ability of a component to be
// injected when method of Rest Controller is requested.
//
//	type Example struct {
//	  at.Scope `value:"request"`
//	  ...
//	}
//
// Deprecated: use Scope instead
type ContextAware struct {
	Annotation

	Scope `value:"request"`
}
