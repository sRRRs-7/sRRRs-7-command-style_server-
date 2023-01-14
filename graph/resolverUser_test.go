package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sRRRs-7/loose_style.git/utils"
	"github.com/stretchr/testify/require"
)

func TestUserResolver(t *testing.T) {
	CreateUserTest(t)
	LoginUserTest(t)
}

func CreateUserTest(t *testing.T) {
	r := GinTestRouter()

	username := utils.RandomString(10)
	email := utils.RandomEmail()
	query := fmt.Sprintf(`
		mutation {
			createUser(username: %s, password: %s, email: %s, sex: %s, date_of_birth: %s) {
				is_error
				message
			}
	}`, fmt.Sprintf("\"%s\"", username), "\"srrrs\"", fmt.Sprintf("\"%s\"", email), "\"man\"", "\"1996-13-45\"")

	q := struct {
		Query string
	}{
		Query: query,
	}

	body := bytes.Buffer{}
	if err := json.NewEncoder(&body).Encode(&q); err != nil {
		t.Fatal("error encode", err)
	}
	req, _ := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(body.Bytes()))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res struct {
		Data struct {
			CreateUser struct {
				IsError bool
				Message string
			}
		}
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, res.Data.CreateUser.IsError, false)
	require.Equal(t, res.Data.CreateUser.Message, "CreateUser OK")
}

func LoginUserTest(t *testing.T) {
	r := GinTestRouter()

	query := fmt.Sprintf(`
		mutation {
			loginUser(username: %s, password: %s) {
				user_id
				username
				OK
			}
	}`, "\"x\"", "\"x\"")

	q := struct {
		Query string
	}{
		Query: query,
	}

	body := bytes.Buffer{}
	if err := json.NewEncoder(&body).Encode(&q); err != nil {
		t.Fatal("error encode", err)
	}
	req, _ := http.NewRequest("POST", "http://localhost:8080/query", bytes.NewBuffer(body.Bytes()))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res struct {
		Data struct {
			LoginUser struct {
				UserId   int
				Username string
				OK       bool
			}
		}
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)

	fmt.Println(w.Body)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, res.Data.LoginUser.UserId, 0)
	require.Equal(t, res.Data.LoginUser.Username, "x")
	require.Equal(t, res.Data.LoginUser.OK, true)
}
