package validate

import "errors"

func Login(email, password string) error {
	if err := UserEmail(email); err != nil {
		return err
	}
	if len(password) < 8 {
		return errors.New("password is too short")
	}
	return nil
}

func Register(email, password string, uid int64) error {
	if err := UserID(uid); err != nil {
		return err
	}
	if err := UserEmail(email); err != nil {
		return err
	}
	if len(password) < 8 {
		return errors.New("password is too short")
	}
	return nil
}
