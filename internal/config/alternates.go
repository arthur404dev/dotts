package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type AlternateMatch struct {
	Path     string
	Score    int
	Suffixes []string
}

type AlternateResolver struct {
	os       string
	distro   string
	profile  string
	hostname string
}

func NewAlternateResolver(osName, distro, profile, hostname string) *AlternateResolver {
	return &AlternateResolver{
		os:       osName,
		distro:   distro,
		profile:  profile,
		hostname: hostname,
	}
}

func (r *AlternateResolver) ResolveFile(basePath string) (string, error) {
	dir := filepath.Dir(basePath)
	base := filepath.Base(basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		if _, statErr := os.Stat(basePath); statErr == nil {
			return basePath, nil
		}
		return "", err
	}

	var candidates []AlternateMatch

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if name == base {
			candidates = append(candidates, AlternateMatch{
				Path:  filepath.Join(dir, name),
				Score: 0,
			})
			continue
		}

		if !strings.HasPrefix(name, base+"##") {
			continue
		}

		suffix := strings.TrimPrefix(name, base+"##")
		suffixes := strings.Split(suffix, ",")

		score, matches := r.calculateScore(suffixes)
		if matches {
			candidates = append(candidates, AlternateMatch{
				Path:     filepath.Join(dir, name),
				Score:    score,
				Suffixes: suffixes,
			})
		}
	}

	if len(candidates) == 0 {
		return basePath, nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	return candidates[0].Path, nil
}

func (r *AlternateResolver) calculateScore(suffixes []string) (int, bool) {
	score := 0

	for _, suffix := range suffixes {
		parts := strings.SplitN(suffix, ".", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := strings.ToLower(parts[1])

		switch key {
		case "hostname":
			if value == strings.ToLower(r.hostname) {
				score += 1000
			} else {
				return 0, false
			}
		case "profile":
			if value == strings.ToLower(r.profile) {
				score += 100
			} else {
				return 0, false
			}
		case "distro":
			if value == strings.ToLower(r.distro) {
				score += 50
			} else {
				return 0, false
			}
		case "os":
			if value == strings.ToLower(r.os) {
				score += 10
			} else {
				return 0, false
			}
		case "default":
			score += 1
		}
	}

	return score, true
}

func (r *AlternateResolver) ResolveDirectory(dir string) (map[string]string, error) {
	result := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		name := info.Name()
		if !strings.Contains(name, "##") {
			return nil
		}

		baseName := strings.Split(name, "##")[0]
		basePath := filepath.Join(filepath.Dir(path), baseName)

		resolved, err := r.ResolveFile(basePath)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(dir, basePath)
		result[relPath] = resolved

		return nil
	})

	return result, err
}

func (r *AlternateResolver) GetAlternatesForFile(basePath string) ([]AlternateMatch, error) {
	dir := filepath.Dir(basePath)
	base := filepath.Base(basePath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var matches []AlternateMatch

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if !strings.HasPrefix(name, base+"##") && name != base {
			continue
		}

		var suffixes []string
		if strings.Contains(name, "##") {
			suffix := strings.TrimPrefix(name, base+"##")
			suffixes = strings.Split(suffix, ",")
		}

		matches = append(matches, AlternateMatch{
			Path:     filepath.Join(dir, name),
			Suffixes: suffixes,
		})
	}

	return matches, nil
}
