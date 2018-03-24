package k8s

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hidevopsio/hi/boot/pkg/log"
	"testing"
	"github.com/stretchr/testify/assert"
)

type Secret struct{
	Name string
	Username string
	Password string
	Namespace string

	secrets *corev1.Secret
}

func init()  {

}

// Create new instance of type Secret
func NewSecret(name, username, password, namespace string) (*Secret) {
	log.Debug("NewSecret")
	s := &Secret{
		Name: name,
		Username: username,
		Password: password,
		Namespace: namespace,
	}

	return s
}

// Create takes the representation of a secret and creates it.  Returns the server's representation of the secret, and an error, if there is any.
func (s *Secret) Create() error  {
	log.Debug("Secret.Create()")
	var data map[string][]byte
	if s.Username != "" {
		data = map[string][]byte{
			"username": []byte(s.Username),
			"password": []byte(s.Password),
		}
	} else {
		data = map[string][]byte{
			"password": []byte(s.Password),
		}
	}

	coreSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Name,
			Labels: map[string]string{
				"username": s.Username,
			},
		},
		Data: data,
	}
	var err error
	s.secrets, err = ClientSet.CoreV1().Secrets(s.Namespace).Create(coreSecret)
	return err
}

func (s *Secret) Get() (*corev1.Secret, error)  {
	log.Debug("Secret.Get()")
	var err error
	s.secrets, err = ClientSet.CoreV1().Secrets(s.Namespace).Get(s.Name, metav1.GetOptions{})

	return s.secrets, err
}

// /Users/johnd/go/src/k8s.io/apimachinery/pkg/apis/meta/v1/types.go
func TestSecretCrud(t *testing.T) {
	secret := NewSecret("test-secret", "test", "tE5t1100", "openshift")

	err := secret.Create()
	assert.Equal(t, nil, err)
}