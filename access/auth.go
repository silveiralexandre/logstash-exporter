package access

import (
	"log"
	"os"
)

// Credential contains all required fields for connecting to the target service
type Credential struct {
	Host     string
	Username string
	Password string
	Limit    int
	Timeout  int
	Retries  int
}

// PullEnvTower fetches from pre-defined environment variables all instance information required to connect
// to a targe Ansible Tower instance
func (a *Credential) PullEnvTower(limit, timeout, retries *int) (Credential, error) {
	const (
		envHost     = "TOWER_HOST"
		envUser     = "TOWER_USER"
		envPassword = "TOWER_PASSWORD"
	)

	// get's values from environment variables
	a.Host = os.Getenv(envHost)
	a.Username = os.Getenv(envUser)
	a.Password = os.Getenv(envPassword)

	// get values from flags or default values
	a.Limit = *limit
	a.Timeout = *timeout
	a.Retries = *retries

	if a.Host == "" || a.Username == "" || a.Password == "" {
		log.Fatal("Required environment variables are not set: $TOWER_HOST, $TOWER_USER, $TOWER_PASSWORD")
	}

	return *a, nil
}
