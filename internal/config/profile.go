package config

import "fmt"

// Profile represents a named sync configuration targeting a specific
// Vault path and local .env output file.
type Profile struct {
	Name      string `mapstructure:"name"`
	VaultPath string `mapstructure:"vault_path"`
	OutputFile string `mapstructure:"output_file"`
	MountPath string `mapstructure:"mount_path"`
}

// Validate checks that all required fields on a Profile are populated.
func (p *Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("profile name must not be empty")
	}
	if p.VaultPath == "" {
		return fmt.Errorf("profile %q: vault_path must not be empty", p.Name)
	}
	if p.OutputFile == "" {
		return fmt.Errorf("profile %q: output_file must not be empty", p.Name)
	}
	return nil
}

// GetProfile returns the Profile with the given name from the Config.
// It returns an error if no profile with that name exists.
func (c *Config) GetProfile(name string) (*Profile, error) {
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return &c.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("profile %q not found", name)
}

// DefaultMountPath returns the mount path for the profile, falling back
// to "secret" if none is explicitly configured.
func (p *Profile) DefaultMountPath() string {
	if p.MountPath == "" {
		return "secret"
	}
	return p.MountPath
}
