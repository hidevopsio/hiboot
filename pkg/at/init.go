package at

// BeforeInit annotation in hiboot is used to init a func / method before the application is initialized
//
//	type Example struct {
//	  at.PostInit
//	  ...
//	}
//

type BeforeInit struct {
	Annotation `json:"-"`

	BaseAnnotation
}

// AfterInit annotation in hiboot is used to init a func / method after the application is initialized
//
//	type Example struct {
//	  at.AfterInit
//	  ...
//	}
//

type AfterInit struct {
	Annotation `json:"-"`

	BaseAnnotation
}
