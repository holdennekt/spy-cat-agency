package catapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CatValidator interface {
	ValidateBreed(breed string) (bool, error)
}

type catValidator struct {
	client *http.Client
}

type CatAPIBreed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCatValidator() CatValidator {
	return &catValidator{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (v *catValidator) ValidateBreed(breed string) (bool, error) {
	url := "https://api.thecatapi.com/v1/breeds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := v.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("cat API returned status %d", resp.StatusCode)
	}

	var breeds []CatAPIBreed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false, err
	}

	breedLower := strings.ToLower(breed)
	for _, b := range breeds {
		if strings.ToLower(b.Name) == breedLower {
			return true, nil
		}
	}

	return false, nil
}
