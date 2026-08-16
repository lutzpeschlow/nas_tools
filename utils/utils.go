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

// ----------------------------------------------------------------------------
//
//	PrintIdArray
//
// ----------------------------------------------------------------------------
func PrintArray(id_array []int) {
	len_id_array := len(id_array)
	// print total array - smaller than 100 entries
	if len_id_array <= 100 {
		for i, id := range id_array {
			if i > 0 && i%10 == 0 {
				fmt.Println()
			}
			fmt.Printf("%8d ", id)
		}
		// large array - print first 50 and last 50 entries
	} else {
		for i := 0; i < 50; i++ {
			if i > 0 && i%10 == 0 {
				fmt.Println()
			}
			fmt.Printf("%8d ", id_array[i])
		}
		fmt.Println("\n...")
		for i := len_id_array - 50; i < len_id_array; i++ {
			if (i-(len_id_array-50)) > 0 && (i-(len_id_array-50))%10 == 0 {
				fmt.Println()
			}
			fmt.Printf("%8d ", id_array[i])
		}
	}
	fmt.Println("\nlength of id array: ", len_id_array)
}
