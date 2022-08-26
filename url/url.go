package url

import (
	"math/rand"
	"net/url"
	"time"
)

const (
	size    = 5
	symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_-+"
)

type Repository interface {
	HasId(id string) bool
	FindById(id string) *Url
	FindByUrl(url string) *Url
	Save(url Url) error
	RegisterClick(id string)
	FindByClick(id string) int
}

type Url struct {
	Id      string    `json:"id"`
	Created time.Time `json:"Created"`
	Final   string    `json:"Final"`
}

type Stats struct {
	Url    *Url `json:"url"`
	Clicks int  `json:"clicks"`
}

var repo Repository

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ConfigRepository(r Repository) {
	repo = r
}

func RegisterClick(id string) {
	repo.RegisterClick(id)
}

func FindOrCreateNewUrl(destiny string) (u *Url, new bool, err error) {
	if u = repo.FindByUrl(destiny); u != nil {
		return u, false, nil
	}

	if _, err = url.ParseRequestURI(destiny); err != nil {
		return nil, false, err
	}

	url := Url{generateId(), time.Now(), destiny}
	repo.Save(url)
	return &url, true, nil
}

func Find(id string) *Url {
	return repo.FindById(id)
}

func (u *Url) Stats() *Stats {
	clicks := repo.FindByClick(u.Id)
	return &Stats{u, clicks}
}

func generateId() string {
	newId := func() string {
		id := make([]byte, size, size)
		for i := range id {
			id[i] = symbols[rand.Intn(len(symbols))]
		}
		return string(id)
	}

	for {
		if id := newId(); !repo.HasId(id) {
			return id
		}
	}
}
