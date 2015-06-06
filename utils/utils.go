package utils

import (
	"fmt"
)

// FIXME
// By dafault, we assume the name is legal.
func ValidateName(name string) error {
	return nil
}

// Duplicate data volumes from source container to target container.
func DuplicateDataVolumes(sVols map[string]string, tVols map[string]string) error {
	for sInPath, sOutPath := range sVols {
		// Traverse source container's data volumes
		// And get the corresponding target container's data volume path
		if tOutPath, ok := tVols[sInPath]; ok {
			// Start to duplicate source data volume's detais to the target
			if err := copyVolume(sOutPath, tOutPath); err != nil {
				fmt.Println("Got an error when duplicating data volume %s to %s", sOutPath, tOutPath)
				return err
			}
		}
	}
	return nil
}

func RemoveCompressedVolumes(id string, index int) error {
	return nil
}
