package config

import "fmt"

// Profile represents a single sync profile configuration.
type Profile struct {
	Name       string `mapstructure:"name"`
	VaultPath  string `mapstructure:"vault_path"`
	OutputFile string `mapstructure:"output_file"`
	MountPath  string `mapstructure:"mount_path"`
	Merge      bool   `mapstructure:"merge"`
}

// Validate checks that the profile has all required fields.
func (p *Profile) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("profile name is required")
	}
	if p.VaultPath == "" {
		return fmt.Errorf("profile %q: vault_path is required", p.Name)
	}
	if p.OutputFile == "" {
		return fmt.Errorf("profile %q: output_file is required", p.Name)
	}
	return nil
}

// DefaultMountPath returns the mount path, falling back to "secret" if empty.
func (p *Profile) DefaultMountPath() string {
	if p.MountPath == "" {
		return "secret"
	}
	return p.MountPath
}

// GetProfile looks up a profile by name from the config.
func (c *Config) GetProfile(name string) (*Profile, error) {
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return &c.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("profile %q not found", name)
}

// GetDefaultProfile returns the first profile if exactly one is defined.
func (c *Config) GetDefaultProfile() (*Profile, error) {
	if len(c.Profiles) == 1 {
		return &c.Profiles[0], nil
	}
	return nil, fmt.Errorf("no default profile: %d profiles defined, specify one explicitly", len(c.Profiles))
}
