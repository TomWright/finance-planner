package repository

import (
	"encoding/json"
	"fmt"
	"github.com/tomwright/finance-planner/internal/application/domain"
	"github.com/tomwright/finance-planner/internal/errs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Profile allows you to load and save a full profile.
type Profile interface {
	// LoadProfile loads the given profile.
	LoadProfile(name string) (*domain.Profile, errs.Error)
	// SaveProfile saves the given profile.
	SaveProfile(profile *domain.Profile) errs.Error
}

func NewProfile(storageDir string) Profile {
	return &jsonProfile{
		storageDir: storageDir,
	}
}

type jsonStoredProfile struct {
	Name string `json:"name"`
}

// jsonProfile implements Profile
type jsonProfile struct {
	storageDir string
}

func (x *jsonProfile) getProfileFilePath(profileName string) string {
	return fmt.Sprintf("%s.json", filepath.Join(x.storageDir, profileName))
}

// LoadProfile loads the given profile.
func (x *jsonProfile) LoadProfile(name string) (*domain.Profile, errs.Error) {
	bytes, err := ioutil.ReadFile(x.getProfileFilePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			// the profile save file did not exist
			return nil, nil
		}
		return nil, errs.FromErr(err).WithCode(errs.ErrCouldNotReadSaveFile)
	}

	var p jsonStoredProfile

	if err := json.Unmarshal(bytes, &p); err != nil {
		return nil, errs.FromErr(err)
	}

	profile := domain.NewProfile()
	profile.Name = p.Name

	return profile, nil
}

// SaveProfile saves the given profile.
func (x *jsonProfile) SaveProfile(profile *domain.Profile) errs.Error {
	storedProfile := jsonStoredProfile{
		Name: profile.Name,
	}

	bytes, err := json.Marshal(storedProfile)
	if err != nil {
		return errs.FromErr(err)
	}

	if err := ioutil.WriteFile(x.getProfileFilePath(profile.Name), bytes, 0644); err != nil {
		return errs.FromErr(err).WithCode(errs.ErrCouldNotWriteSaveFile)
	}

	return nil
}
