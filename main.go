package main

import (
	"encoding/json"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"net/http"
	"os"
	"strings"
)

/**
 * @Author: yangyang
 * @Date: 2022/12/8 17:28
 * @Desc:
 */

func main() {
	app := kingpin.New("jwks", "envoy remote jwks go.")
	jwk := app.Command("jwks", "get jwks by file web server")
	var configFile string
	var port int
	var address string
	var b []byte
	parseConfig := func(_ *kingpin.ParseContext) error {

		f, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer f.Close()
		b, err = io.ReadAll(f)
		if err != nil {
			return err
		}
		return nil
	}
	jwk.Flag("file-path", "jwks file path").Short('c').PlaceHolder("/path/to/file").Action(parseConfig).ExistingFileVar(&configFile)
	jwk.Flag("port", "http port").Short('p').Default("8080").PlaceHolder("<port>").IntVar(&port)
	jwk.Flag("address", "http address").Short('d').Default("0.0.0.0").PlaceHolder("<ipaddr>").StringVar(&address)
	app.HelpFlag.Short('h')
	args := os.Args[1:]
	cmd := kingpin.MustParse(app.Parse(args))
	switch cmd {
	case jwk.FullCommand():
		if configFile == "" {
			fmt.Println("config file is must")
			os.Exit(2)
		}
		http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("success"))
		})
		http.HandleFunc("/jwks", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write(b)
		})
		http.HandleFunc("/valid", func(writer http.ResponseWriter, request *http.Request) {
			message := json.RawMessage(b)
			jwks, err := keyfunc.NewJSON(message)
			if err != nil {
				writer.Write([]byte(fmt.Sprintf("jwks valid fail: %s", err.Error())))
				return
			}
			var token string
			reqToken := request.Header.Get("Authorization")
			if reqToken != "" {
				splitToken := strings.Split(reqToken, "Bearer ")
				if len(splitToken) == 2 {
					token = splitToken[1]
				}
			}
			if token == "" {
				token = request.URL.Query().Get("access_token")
			}
			if token == "" {
				writer.Write([]byte("token is empty"))
				return
			}
			parse, err := jwt.Parse(token, jwks.Keyfunc)
			if err != nil {
				writer.Write([]byte(fmt.Sprintf("Failed to parse the JWT.\nError: %s", err.Error())))
				return
			}
			if !parse.Valid {
				writer.Write([]byte("The token is not valid."))
				return

			}
			writer.Write([]byte("The token is valid."))
		})
		fmt.Println(fmt.Sprintf("listen address is :%s:%d", address, port))
		http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)

	default:
		app.Usage(args)
		os.Exit(2)
	}
}
