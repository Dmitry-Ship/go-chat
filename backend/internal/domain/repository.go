package domain

type GenericRepository[T Aggregate] interface {
	Store(aggregate T) error
	Update(aggregate T) error
}
