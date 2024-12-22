package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/zalhui/calc_golang/pkg/calculation"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")

	if config.Addr == "" {
		config.Addr = "8080"
	}

	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		request := new(Request)
		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(request)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "Bad request"}`)
			return
		}

		result, err := calculation.Calc(request.Expression)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			responce := Response{Error: err.Error()}

			json.NewEncoder(w).Encode(responce)
		} else {
			responce := Response{Result: result}
			json.NewEncoder(w).Encode(responce)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error": "only POST method allowed"}`)
	}

}

func (a *Application) RunServer() error {
	http.HandleFunc("/", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
