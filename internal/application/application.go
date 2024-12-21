package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/zalhui/calc_golang/pkg/calculation"
)

type Config struct{
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

type Request struct{
	Expression string `json:"expression`
}
 
func CalcHandler(w http.ResponseWriter, r *http.Request){
	request:=new(Request)
	defer r.Body.Close()

	err:=json.NewDecoder(r.Body).Decode(request)

	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest
		return 
	}

	result,err:=calculation.Calc(request.Expression)

 

		if err != nil {
			errorMessages := map[error]string{
				calculation.ErrInvalidExpr:    "err:%s",
				calculation.ErrEmptyExpression: "err:%s",
				calculation.ErrDivisionByZero: "err:%s",
			}
		
			if msg, ok := errorMessages[err]; ok {
				fmt.Fprintf(w, msg, err.Error())
			} else {
				fmt.Fprintf(w, "unknown err")
			}
		}
		
}

func (a *Application) RunServer() error{
	http.HandleFunc("/", CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}