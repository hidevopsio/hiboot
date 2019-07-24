package at

// RequestBody the annotation RequestBody
type RequestBody interface{}

// RequestForm the annotation RequestForm
// TODO: should investigate why interface is not working on this annotation
// TODO: see uint test "should return http.StatusInternalServerError when input form field validation failed"
type RequestForm struct{}

// RequestParams the annotation RequestParams
type RequestParams interface{}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty" json:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty" json:"per_page,omitempty"`
}