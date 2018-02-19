package models

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type Fishpic struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Species   string    `json:"species" db:"species"`
	Location  string    `json:"location" db:"location"`
}

// String is not required by pop and may be deleted
func (f Fishpic) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

// Fishpics is not required by pop and may be deleted
type Fishpics []Fishpic

// String is not required by pop and may be deleted
func (f Fishpics) String() string {
	jf, _ := json.Marshal(f)
	return string(jf)
}

func (f *Fishpic) Init() {
	newUUID, err := uuid.NewV4()
	if err != nil {
		errors.WithStack(err)
	}
	f.ID = newUUID

	currTime := time.Now().UTC()
	f.CreatedAt = currTime
	f.UpdatedAt = currTime
}

func (f *Fishpic) CreateAndSave() error {
	f.Init()
	err := f.ValidateAndSave()
	if err != nil {
		return err
	}

	return nil
}

func (f *Fishpic) ValidateAndSave() error {
	bytes, err := json.Marshal(f)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./db/"+f.ID.String(), bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}
