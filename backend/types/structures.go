package types

type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn string
	}
	Jwt struct {
		Secret string
	}
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}
