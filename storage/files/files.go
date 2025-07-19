package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"

	"notekeeper-tgbot/lib/e"
	"notekeeper-tgbot/storage"
)

const defaultPerm = 7777

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can`t save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can`t pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return nil, storage.ErrNoSavedPages
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := randRange(0, len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		msg := fmt.Sprintf("can`t check file %s exist", fileName)
		return e.Wrap(msg, err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can`t remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		msg := fmt.Sprintf("can`t check file %s exist", fileName)
		return false, e.Wrap(msg, err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can`t check file %s exist", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.WrapIfErr("can`t open page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.WrapIfErr("can`t decode page", err)
	}

	return &p, nil
}
