package model

type Lists interface {
	HasMember(id string) bool
	Delete(id string) error
	Add(id, name string) error
	Update(id, name string) error
}
