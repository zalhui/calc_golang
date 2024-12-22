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
	Result string `json:"result"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Internal server error"}`)
		return
	}

	result, err := calculation.Calc(request.Expression)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err == calculation.ErrInvalidExpression || err == calculation.ErrBrackets {
			fmt.Fprintf(w, `{"error": "Expression is not valid"}`)
		} else {
			fmt.Fprintf(w, `{"error": "unknown err"}`)
		}
	} else {
		json.NewEncoder(w).Encode(Response{Result: fmt.Sprintf("%f", result)})
	}

}

func (a *Application) RunServer() error {
	http.HandleFunc("/", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
