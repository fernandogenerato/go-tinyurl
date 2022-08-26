package url

type repositoryMemory struct {
	urls   map[string]*Url
	clicks map[string]int
}

func (r *repositoryMemory) FindByClick(id string) int {
	return r.clicks[id]
}

func NewMemoryRepository() *repositoryMemory {
	return &repositoryMemory{
		make(map[string]*Url),
		make(map[string]int),
	}
}

func (r *repositoryMemory) HasId(id string) bool {
	_, has := r.urls[id]
	return has
}

func (r *repositoryMemory) FindById(id string) *Url {
	return r.urls[id]
}

func (r *repositoryMemory) FindByUrl(url string) *Url {
	for _, u := range r.urls {
		if u.Final == url {
			return u
		}
	}
	return nil
}

func (r *repositoryMemory) Save(url Url) error {
	r.urls[url.Id] = &url
	return nil
}

func (r *repositoryMemory) RegisterClick(id string) {
	r.clicks[id] += 1
}

func (r *repositoryMemory) BuscarClicks(id string) int {
	return r.clicks[id]
}
