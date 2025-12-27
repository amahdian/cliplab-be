package env

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
)

// Envs contains all application configurations.
// The configuration are parsed from environment variables.
// For more info see: https://github.com/sethvargo/go-envconfig
type Envs struct {
	Server struct {
		GinMode         string `env:"GIN_MODE, default=debug"`
		LogLevel        string `env:"LOG_LEVEL, default=debug"`
		LogFormat       string `env:"LOG_FORMAT, default=text"`
		HttpPort        string `env:"HTTP_PORT, default=8080"`
		SwaggerHostAddr string `env:"SWAGGER_HOST_ADDR"`
		LocalCors       bool   `env:"LOCAL_CORS, default=false"`
		AssetsDir       string `env:"ASSETS_DIR, default=/assets"`
		JwtSecret       string `env:"JWT_SECRET, required"`
		ApiHostAddr     string `env:"API_HOST, required"`
	}

	Db struct {
		LogLevel string `env:"LOG_LEVEL, default=error"`
		Dsn      string `env:"DSN, required"`
	}

	Gemini struct {
		ClientHost string `env:"GEMINI_HOST, required"`
		Token      string `env:"GEMINI_TOKEN, required"`
	}

	RapidApi struct {
		Token string `env:"RAPIDAPI_TOKEN, required"`
	}

	Redis struct {
		Address  string `env:"REDIS_ADDRESS, default=localhost:6379"`
		Password string `env:"REDIS_PASSWORD"`
		DB       int    `env:"REDIS_DB, default=0"`
	}

	FileStorage struct {
		Bypass    bool   `env:"BYPASS_STORAGE, default=true"`
		Endpoint  string `env:"S3_ENDPOINT, default=http://localhost:9090"`
		AccessKey string `env:"S3_ACCESS_KEY"`
		SecretKey string `env:"S3_SECRET_KEY"`
		Region    string `env:"S3_REGION"`
	}
}

// Load loads the environment variables from the .env files
// we use a base env file called ".env" and override it for different environment (e.g. dev, test, prod, or testing)
func Load(basePath string) (*Envs, error) {
	envFiles, err := filepath.Glob(filepath.Join(basePath, ".env*"))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find env files.")
	}

	profileToEnvFile := map[string]string{}
	for _, envFile := range envFiles {
		// get the profile from the env file name e.g. .env.dev -> dev
		profile := filepath.Base(envFile)
		profile, _ = strings.CutPrefix(profile, ".env")
		profile, _ = strings.CutPrefix(profile, ".")
		profileToEnvFile[profile] = envFile
	}

	// first let load the env file for the current profile
	profile := os.Getenv("PROFILE")
	activeEnvFiles := make([]string, 0)
	if profile != "" {
		envFile, ok := profileToEnvFile[profile]
		if !ok {
			return nil, fmt.Errorf("No env file found for %q profile.", profile)
		}
		activeEnvFiles = append(activeEnvFiles, envFile)
	}

	// then let's add the default profile as well if it exists
	if envFile, ok := profileToEnvFile[""]; ok {
		activeEnvFiles = append(activeEnvFiles, envFile)
	}

	if len(activeEnvFiles) == 0 {
		log.Print("no .env file found")
	}

	log.Printf(fmt.Sprintf("trying to load %v env files", activeEnvFiles))
	err = godotenv.Load(activeEnvFiles...)
	if err != nil {
		return nil, err
	}

	var configs Envs
	if err = envconfig.Process(context.Background(), &configs); err != nil {
		return nil, err
	}

	return &configs, nil
}
