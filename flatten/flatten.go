package flatten

import (
	"errors"
	"fmt"
	"os"
	"path"
)

// Flatten a directory recursively
func Flatten(name string) error {
	dir, err := os.ReadDir(name)
	if err != nil {
		return fmt.Errorf("listing directory: %w", err)
	}

	for _, item := range dir {
		if !item.IsDir() {
			continue
		}

		err := flattenSubdir(name, path.Join(name, item.Name()))
		if err != nil {
			return fmt.Errorf("flattening subdirectory: %w", err)
		}
	}

	return nil
}

func flattenSubdir(flattenTarget, name string) error {
	dir, err := os.ReadDir(name)
	if err != nil {
		return fmt.Errorf("listing directory: %w", err)
	}

	for _, sub := range dir {
		if sub.IsDir() {
			err := flattenSubdir(flattenTarget, path.Join(name, sub.Name()))
			if err != nil {
				return err
			}
			continue
		}

		oldPath := path.Join(name, sub.Name())
		newPath := path.Join(flattenTarget, sub.Name())

		_, err := os.Stat(newPath)
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("flattening %s: %s: %w", oldPath, newPath, os.ErrExist)
		}

		err = os.Rename(oldPath, newPath)
		if err != nil {
			return fmt.Errorf("moving file: %w", err)
		}
	}

	if err := os.Remove(name); err != nil {
		return fmt.Errorf("removing emptied directory: %w", err)
	}
	return nil
}
