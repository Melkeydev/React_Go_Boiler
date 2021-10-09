package types 

type Config struct {
  Port int
  Env string
  Db struct {
    Dsn string
  }
}
