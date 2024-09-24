package user

import (
	"encoding/json"
	"net/http"

	user_dal "github.com/Arinji2/meme-backend/sql/dal/user"
	"github.com/Arinji2/meme-backend/types"
	"github.com/go-chi/render"
)

func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req types.EmailToUserInput
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}
	user, err := user_dal.GetUserByEmail(r.Context(), req.Email)

	if err != nil {

		render.Status(r, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}
