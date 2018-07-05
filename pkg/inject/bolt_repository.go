package inject

type BoltRepository struct{
	Repository
	Bucket string `json:"bucket"`
}
