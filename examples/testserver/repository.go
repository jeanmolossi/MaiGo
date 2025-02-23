package testserver

var _ Repository = (*Impl)(nil)

type (
	Repository interface {
		GetAll() (models []Model)
		GetByID(id uint) (model Model, found bool)
		Create(model Model) (created Model)
		Delete(id uint)
	}

	Impl struct {
		table *Table
	}

	Provider struct {
		User     Repository
		Resource Repository
	}
)

// Create implements Repository.
func (i *Impl) Create(model Model) (created Model) {
	newID := uint(len(*i.table) + 1)
	model.SetID(newID)
	(*i.table)[newID] = model

	return model
}

// Delete implements Repository.
func (i *Impl) Delete(id uint) {
	delete(*i.table, id)
}

// GetAll implements Repository.
func (i *Impl) GetAll() (models []Model) {
	for _, model := range *i.table {
		models = append(models, model)
	}

	return
}

// GetByID implements Repository.
func (i *Impl) GetByID(id uint) (model Model, found bool) {
	model, found = (*i.table)[id]

	return
}

func newRepositoryImpl(table *Table) *Impl {
	return &Impl{
		table: table,
	}
}

func NewProvider(state *State) *Provider {
	userNamespace := (*state)[UserNamespace]
	resourceNamespace := (*state)[ResourceNamespace]

	return &Provider{
		User:     newRepositoryImpl(&userNamespace),
		Resource: newRepositoryImpl(&resourceNamespace),
	}
}
