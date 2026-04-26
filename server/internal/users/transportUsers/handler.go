package transportUsers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"lab2pepe/server/internal/core/domain"
	"lab2pepe/server/internal/core/httpUtils"
	"lab2pepe/server/internal/users/serviceUsers"
	"math/big"
	"net/http"
)

type HandlersUsers struct {
	ServiceUsers *serviceUsers.ServiceUsers
}

func NewHandlersUser(service *serviceUsers.ServiceUsers) *HandlersUsers {
	return &HandlersUsers{ServiceUsers: service}
}

func (h *HandlersUsers) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var dto AddUsersDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := dto.ValidateForCreateUser(); err != nil {
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}
	if err := h.ServiceUsers.AddUsers(r.Context(), dto.Login, dto.Password); err != nil {
		httpUtils.WriteError(w, err, http.StatusConflict)
		return
	}
	b, err := json.MarshalIndent(CreateUserResponseDTO{Login: dto.Login}, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
	}
}

func (h *HandlersUsers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var dto LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	result, err := h.ServiceUsers.Login(r.Context(), dto.Login, dto.Password)
	if err != nil {
		httpUtils.WriteError(w, err, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponseDTO{
		Token:      result.Token,
		TmpToken:   result.TmpToken,
		Require2FA: result.Require2FA,
	})
}

func (h *HandlersUsers) HandleLogin2FA(w http.ResponseWriter, r *http.Request) {
	var dto Login2FADTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	token, err := h.ServiceUsers.LoginConfirm2FA(r.Context(), dto.TmpToken, dto.Code)
	if err != nil {
		fmt.Println("HandleLogin2FA error:", err)
		httpUtils.WriteError(w, err, http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponseDTO{Token: token})
}

func (h *HandlersUsers) HandleSetup2FA(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDKey).(int)
	if !ok {
		httpUtils.WriteError(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}
	login, err := h.ServiceUsers.GetLoginByID(r.Context(), userID)
	if err != nil {
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	otpauthURL, secret, err := h.ServiceUsers.Setup2FA(r.Context(), userID, login)
	if err != nil {
		fmt.Println("HandleSetup2FA error:", err)
		httpUtils.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Setup2FAResponseDTO{
		OtpauthURL: otpauthURL,
		Secret:     secret,
	})
}

func (h *HandlersUsers) HandleConfirm2FA(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDKey).(int)
	if !ok {
		httpUtils.WriteError(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}
	var dto Confirm2FADTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}
	if err := h.ServiceUsers.Confirm2FA(r.Context(), userID, dto.Code); err != nil {
		fmt.Println("HandleConfirm2FA error:", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"2fa enabled"}`))
}

func (h *HandlersUsers) HandleDelete2FA(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(domain.UserIDKey).(int)
	if !ok {
		httpUtils.WriteError(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}
	var dto Delete2FADTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}
	if err := h.ServiceUsers.Delete2FA(r.Context(), userID, dto.Code); err != nil {
		fmt.Println("HandleDelete2FA error:", err)
		httpUtils.WriteError(w, err, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"2fa disabled"}`))
}

func (h *HandlersUsers) HandleKeyExchange(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(domain.UserIDKey).(int) // из JWT middleware

	var req keyExchangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	nBytes, err := hex.DecodeString(req.N)
	if err != nil {
		http.Error(w, "invalid n", http.StatusBadRequest)
		return
	}
	eBytes, err := hex.DecodeString(req.E)
	if err != nil {
		http.Error(w, "invalid e", http.StatusBadRequest)
		return
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	// Генерируем случайный AES-128 ключ
	aesKey := make([]byte, 16)
	if _, err := rand.Read(aesKey); err != nil {
		http.Error(w, "key generation error", http.StatusInternalServerError)
		return
	}

	// Шифруем AES-ключ через RSA клиента: C = AES^e mod n
	aesInt := new(big.Int).SetBytes(aesKey)
	encrypted := new(big.Int).Exp(aesInt, e, n)

	// Сохраняем AES-ключ за пользователем
	if err := h.ServiceUsers.SaveAESKey(r.Context(), userID, aesKey); err != nil {
		http.Error(w, "failed to save key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyExchangeResponse{
		EncryptedAESKey: hex.EncodeToString(encrypted.Bytes()),
	})
}
