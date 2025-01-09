package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/marcoaureliojf/streamStudio/backend/internal/middlewares"
	"github.com/pion/webrtc/v3"
)

type SignalingHandler struct {
}

func NewSignalingHandler() *SignalingHandler {
	return &SignalingHandler{}
}

type SDPOfferRequest struct {
	SDP string `json:"sdp"`
}
type SDPOfferResponse struct {
	SDP string `json:"sdp"`
}

func (h *SignalingHandler) Offer(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetUserFromContext(r.Context())
	if user == nil {
		log.Println("Usuário não autenticado")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Não autorizado"})
		return
	}

	var request SDPOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Erro ao decodificar o corpo da requisição:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Corpo da requisição inválido"})
		return
	}

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  request.SDP,
	}

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		log.Println("Erro ao setar a descrição remota", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar a requisição"})
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Println("Erro ao criar o answer", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar a requisição"})
		return
	}
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		log.Println("Erro ao setar a descrição local", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Erro ao processar a requisição"})
		return
	}

	response := SDPOfferResponse{
		SDP: answer.SDP,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
