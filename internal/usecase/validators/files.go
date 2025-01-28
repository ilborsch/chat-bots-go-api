package validate

import "fmt"

func File(id, ownerID int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid file id provided %v", id)
	}
	if err := UserID(ownerID); err != nil {
		return err
	}
	return nil
}

func SaveFile(filename string, fileData []byte) error {
	if filename == "" {
		return fmt.Errorf("empty filename is provided")
	}
	if len(fileData) == 0 {
		return fmt.Errorf("empty file is provided")
	}
	return nil
}
