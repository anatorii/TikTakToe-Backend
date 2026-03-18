package datasource

type Postgres struct {
	Pool *DbPool
}

func NewPostgres() Postgres {
	return Postgres{
		Pool: NewDbPool(NewPoolConfig("")),
	}
}
