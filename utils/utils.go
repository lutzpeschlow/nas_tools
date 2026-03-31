package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ----------------------------------------------------------------------------
//
//	WriteTxtFile
//
// ----------------------------------------------------------------------------
func WriteTxtFile(filePath string, lines []string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create output file %s: %w", filePath, err)
	}
	defer out.Close()
	for _, line := range lines {
		if _, err := fmt.Fprintln(out, line); err != nil {
			return fmt.Errorf("write output file %s: %w", filePath, err)
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
//
//	WriteSessionFile
//
// ----------------------------------------------------------------------------
func WriteSessionFile(groupName, entityType, outputPath string, ids []string) error {
	sesPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".ses"
	f, err := os.Create(sesPath)
	if err != nil {
		return fmt.Errorf("create session file %s: %w", sesPath, err)
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "ga_group_create(%q)\n", groupName); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(f, "ga_group_entity_add(%q, %q // @ \n", groupName, entityType); err != nil {
		return err
	}
	for i, id := range ids {
		if i == len(ids)-1 {
			if _, err := fmt.Fprintf(f, "\" %s \" )\n", id); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(f, "\" %s \" // @ \n", id); err != nil {
				return err
			}
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
//
//	RemoveDuplicateEntries
//
// ----------------------------------------------------------------------------
func RemoveDuplicateEntries(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))
	for _, v := range input {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
