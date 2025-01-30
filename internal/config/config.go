package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `env:"ENV" env-required:"true" comment:"Среда выполнения приложения"`
	HTTPServer HTTPServer
	Database   Database
	Limiter    Limiter
	JWT        JWT
}

type HTTPServer struct {
	Port           string        `env:"PORT" env-default:"8080" comment:"Порт для HTTP сервера"`
	Timeout        time.Duration `env:"HTTP_SERVER_TIMEOUT" env-default:"4s" comment:"Таймаут для HTTP сервера"`
	IdleTimeout    time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s" comment:"Таймаут бездействия для HTTP сервера"`
	SwaggerEnabled bool          `env:"HTTP_SERVER_SWAGGER_ENABLED" comment:"Включить Swagger"`
}

type Database struct {
	Net                string        `env:"DATABASE_NET" env-default:"tcp" comment:"Тип сети для подключения к базе данных"`
	Host               string        `env:"DATABASE_HOST" env-required:"true" comment:"Хост базы данных"`
	Port               string        `env:"DATABASE_PORT" env-required:"true" comment:"Порт базы данных"`
	DBName             string        `env:"DATABASE_DB_NAME" env-required:"true" comment:"Имя базы данных"`
	User               string        `env:"DATABASE_USER" env-required:"true" comment:"Пользователь базы данных"`
	Password           string        `env:"DATABASE_PASSWORD" env-required:"true" comment:"Пароль пользователя базы данных"`
	SSLMode            string        `env:"DATABASE_SSLMODE" env-default:"disable" comment:"Режим SSL для подключения к базе данных"`
	TimeZone           string        `env:"DATABASE_TIME_ZONE" env-default:"UTC" comment:"Часовой пояс базы данных"`
	Timeout            time.Duration `env:"DATABASE_TIMEOUT" env-default:"2s" comment:"Таймаут подключения к базе данных"`
	MaxIdleConnections int           `env:"DATABASE_MAX_IDLE_CONNECTIONS" env-default:"40" comment:"Максимальное количество простых соединений"`
	MaxOpenConnections int           `env:"DATABASE_MAX_OPEN_CONNECTIONS" env-default:"40" comment:"Максимальное количество открытых соединений"`
}

type Limiter struct {
	RPS   int           `env:"LIMITER_RPS" env-default:"10" comment:"Запросов в секунду"`
	Burst int           `env:"LIMITER_BURST" env-default:"20" comment:"Максимальное количество запросов в секунду"`
	TTL   time.Duration `env:"LIMITER_TTL" env-default:"10m" comment:"Время жизни лимита"`
}

type JWT struct {
	SecretKey       string        `env:"JWT_SECRET_KEY" env-default:"notasecret" comment:"Секретный ключ для JWT"`
	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" env-default:"15m" comment:"Время жизни access токена"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" env-default:"720h" comment:"Время жизни refresh токена"`
}

func MustLoad() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		godotenv.Load()
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Panic(err)
	}

	return &cfg
}
