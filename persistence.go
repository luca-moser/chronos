package chronos

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Storage interface {
	SaveSchedule(schedule TaskSchedule) error
	LoadSchedule() (TaskSchedule, error)
}

func NewFileStorage(path string, filename string) *FileStorage {
	return &FileStorage{Path: path, Filename: filename}
}

type FileStorage struct {
	Path     string
	Filename string
}

func (fs *FileStorage) SaveSchedule(schedule *TaskSchedule) error {

	// ensure path exists
	os.MkdirAll(fs.Path, 0666)

	jsonSchedule, err := json.Marshal(schedule)
	if err != nil {
		return err
	}

	// truncate existing file
	os.Truncate(path.Join(fs.Path, fs.Filename), 0)

	file, err := os.OpenFile(path.Join(fs.Path, fs.Filename), os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(jsonSchedule); err != nil {
		return err
	}

	return nil
}

func (fs *FileStorage) LoadSchedule() (*TaskSchedule, error) {
	fileContent, err := ioutil.ReadFile(path.Join(fs.Path, fs.Filename))
	if err != nil {
		return nil, err
	}

	schedule := &TaskSchedule{}
	if err := json.Unmarshal(fileContent, schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (fs *FileStorage) RemoveSchedule() error {
	return os.Remove(path.Join(fs.Path, fs.Filename))
}
