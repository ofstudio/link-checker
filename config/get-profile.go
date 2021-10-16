package config

import (
	"fmt"
	"os"
)

func MustGetProfile(filename string) *Profile {
	if len(os.Args) != 2 {
		fmt.Printf("Usage\n\n  link-checker PROFILE\n\nSee %s for profiles\n", filename)
		os.Exit(0)
	}

	id := os.Args[1]
	cfg, err := readYaml(filename)
	if err != nil {
		fmt.Printf("Error loading %s: %v\n", filename, err)
		os.Exit(-1)
	}

	if profile, ok := cfg.Profiles[id]; ok {
		profile.Id = id
		return &profile
	}

	fmt.Printf("Unknown profile: %s\n\nSee %s for profiles\n", id, filename)
	os.Exit(-1)
	return nil
}
