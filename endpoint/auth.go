package endpoint

import(
	"net/http"	
	"api"
	"time"
	"auth"
	"fmt"
	"github.com/streadway/simpleuuid"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
)

type AuthEndpoint struct{}

var Auth AuthEndpoint

func (endpoint * AuthEndpoint) IssueToken(r * http.Request, quizApi * api.QuizzicalAPI) (interface{}, error){

	sub := r.FormValue(ParamSubject)

	if len(sub) == 0{
		return nil, fmt.Errorf("Subject Required")
	}

	if sub != "com.asfour.Quizzical" {
		return nil,fmt.Errorf("Subject Invalid")
	}

	iat := time.Now()
	exp := iat.Add(time.Hour)
	jti,err := simpleuuid.NewTime(iat)

	if err != nil {
	 	return nil,err
	 } 

	token :=  jwt.New(jwt.SigningMethodHS256)
	token.Claims["iat"] = iat.Unix()
	token.Claims["nbf"] = token.Claims["iat"]
	token.Claims["exp"] = exp.Unix()
	token.Claims["jti"] = jti
	
	sToken,err := token.SignedString([]byte(auth.Key))

	if err != nil {
		return nil, err
	}

	return api.Response{Data: api.NewToken(sToken,exp)},nil
}