package configs

import (
    "fmt"
    "os"
    "strconv"
    "bufio"
    "strings"
    "syscall"

    "golang.org/x/term"
    "github.com/joho/godotenv"
    
)

type Config struct {
      PublicHost  string
      Port        string
      DBUser      string
      DBPassword  string
      DBAddress   string
      DBName      string
      
}
func credentials() (string, string, error) {
  reader := bufio.NewReader(os.Stdin)

  fmt.Print("Enter Username: ")
  username, err := reader.ReadString('\n')
  if err != nil {
    return "", "", err
  }

  fmt.Print("Enter Password: ")
  bytePassword, err := term.ReadPassword(int(syscall.Stdin))
  if err != nil {
    return "", "", err
  }

  password := string(bytePassword)
  return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

var Envs = initConfig()


func initConfig() Config {
username, password, _ := credentials()

godotenv.Load()

        return Config{
          PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
                Port: getEnv("PORT", "8080"),
                DBUser: getEnv("DB_USER", username),
                DBPassword: getEnv("DB_PASSWORD", password),
                DBAddress: fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
                DBName: getEnv("DB_NAME", "overall_database"),

        }
}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
