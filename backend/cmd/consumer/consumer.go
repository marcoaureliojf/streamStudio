// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/google/uuid"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
// 	"github.com/marcoaureliojf/streamStudio/backend/internal/queue"
// 	"github.com/pion/webrtc/v4"
// 	"github.com/pion/webrtc/v4/pkg/media/ivfwriter"
// 	"github.com/streadway/amqp"
// )

// var peerConnection *webrtc.PeerConnection

// func mainConsumer() {
// 	cfg := config.LoadConfig()
// 	database.Connect(cfg)

// 	rabbitMQ, err := queue.NewRabbitMQ(cfg)
// 	if err != nil {
// 		log.Fatal("Erro ao iniciar RabbitMQ: ", err)
// 	}
// 	defer rabbitMQ.Close()

// 	rabbitMQ.Consume(scheduleHandler())
// }

// func scheduleHandler() func(d amqp.Delivery) {
// 	return func(d amqp.Delivery) {
// 		log.Printf("Recebida mensagem da fila %s\n", d.RoutingKey)

// 		var schedule models.Schedule
// 		err := json.Unmarshal(d.Body, &schedule)
// 		if err != nil {
// 			log.Printf("Erro ao desserializar mensagem: %v", err)
// 			return
// 		}
// 		log.Printf("Processando o agendamento %v", schedule)

// 		stream, err := findStream(schedule.StreamID)
// 		if err != nil {
// 			log.Printf("Erro ao buscar transmissão: %v", err)
// 			return
// 		}

// 		if peerConnection != nil {
// 			log.Println("Já existe uma transmissão em andamento")
// 			return
// 		}

// 		startStream(stream)
// 		log.Printf("Agendamento %v finalizado\n", schedule)
// 		d.Ack(false)
// 	}
// }

// func findStream(id uint) (*models.Stream, error) {
// 	db := database.GetDB()
// 	var stream models.Stream
// 	result := db.First(&stream, id)
// 	if result.Error != nil {
// 		return nil, fmt.Errorf("stream não encontrada")
// 	}
// 	return &stream, nil
// }

// func startStream(stream *models.Stream) {
// 	var err error

// 	// Configuração inicial do PeerConnection
// 	peerConnection, err = webrtc.NewPeerConnection(webrtc.Configuration{
// 		ICEServers: []webrtc.ICEServer{
// 			{
// 				URLs: []string{"stun:stun.l.google.com:19302"},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		log.Fatalf("Erro ao criar PeerConnection: %v", err)
// 	}

// 	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
// 		fmt.Printf("Estado da conexão ICE mudou para %s\n", connectionState.String())
// 		if connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateDisconnected {
// 			peerConnection = nil
// 		}
// 	})

// 	// Configura evento ao receber uma nova faixa de mídia
// 	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
// 		log.Printf("Recebido track: %s \n", track.ID())

// 		id := uuid.New()
// 		fileName := fmt.Sprintf("%s.ivf", id)

// 		// Criação do writer para o arquivo IVF
// 		videoWriter, err := ivfwriter.New(fileName)
// 		if err != nil {
// 			log.Printf("Erro ao criar o writer: %v", err)
// 			return
// 		}
// 		defer func() {
// 			if err = videoWriter.Close(); err != nil {
// 				log.Printf("Erro ao fechar o writer: %v", err)
// 			}
// 		}()

// 		// Loop para receber e salvar pacotes RTP
// 		for {
// 			packet, _, readErr := track.ReadRTP()
// 			if readErr != nil {
// 				log.Printf("Erro ao ler o pacote RTP: %v\n", readErr)
// 				return
// 			}

// 			if writeErr := videoWriter.WriteRTP(packet); writeErr != nil {
// 				log.Printf("Erro ao escrever pacote RTP: %v\n", writeErr)
// 				return
// 			}
// 		}
// 	})

// 	log.Printf("Iniciando a transmissão: %v", stream)
// }

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/marcoaureliojf/streamStudio/backend/internal/config"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models"
	"github.com/marcoaureliojf/streamStudio/backend/internal/queue"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/ivfwriter"
	"github.com/streadway/amqp"
)

var peerConnection *webrtc.PeerConnection
var mediaEngine *webrtc.MediaEngine

func main() {
	cfg := config.LoadConfig()
	database.Connect(cfg)

	queue.Init(cfg)

	rabbitMQ, err := queue.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal("Erro ao iniciar RabbitMQ: ", err)
	}
	defer rabbitMQ.Close()

	mediaEngine = &webrtc.MediaEngine{}
	mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8, ClockRate: 90000, Channels: 0, SDPFmtpLine: ""},
		PayloadType:        100,
	}, webrtc.RTPCodecTypeVideo)

	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))
	rabbitMQ.Consume(scheduleHandler(api))
}

func scheduleHandler(api *webrtc.API) func(d amqp.Delivery) {
	return func(d amqp.Delivery) {
		log.Printf("Recebida mensagem da fila %s\n", d.RoutingKey)

		var schedule models.Schedule
		err := json.Unmarshal(d.Body, &schedule)
		if err != nil {
			log.Printf("Erro ao desserializar mensagem: %v", err)
			return
		}
		log.Printf("Processando o agendamento %v", schedule)

		stream, err := findStream(schedule.StreamID)
		if err != nil {
			log.Printf("Erro ao buscar transmissão: %v", err)
			return
		}

		if peerConnection != nil {
			log.Println("Já existe uma transmissão em andamento")
			return
		}

		startStream(stream, api)

		log.Printf("Agendamento %v finalizado\n", schedule)
		d.Ack(false)
	}
}

func findStream(id uint) (*models.Stream, error) {
	db := database.GetDB()
	var stream models.Stream
	result := db.First(&stream, id)
	if result.Error != nil {
		return nil, fmt.Errorf("stream não encontrada")
	}
	return &stream, nil
}

func startStream(stream *models.Stream, api *webrtc.API) {
	// Configurações iniciais do WebRTC
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		log.Printf("Erro ao criar PeerConnection: %v", err)
		return
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Estado da conexão ICE mudou para %s\n", connectionState.String())
		if connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateDisconnected {
			peerConnection = nil
		}
	})

	// Define o evento que é disparado quando uma nova faixa de mídia é recebida
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Recebido track: %s \n", track.ID())

		buf := new(bytes.Buffer)

		videoWriter, err := ivfwriter.New(buf.String(), nil)
		if err != nil {
			log.Println("Erro ao criar o writer:", err)
			return
		}
		defer func() {
			if err = videoWriter.Close(); err != nil {
				log.Println("Erro ao fechar o writer:", err)
				return
			}
		}()

		//gera um uuid para nome do arquivo
		id := uuid.New()
		// Cria um arquivo para salvar o video
		file, err := os.Create(fmt.Sprintf("%s.ivf", id))
		if err != nil {
			log.Println("Erro ao criar arquivo:", err)
			return
		}
		defer file.Close()

		for {
			// Leitura do pacote RTP
			packet, _, readErr := track.ReadRTP()
			if readErr != nil {
				if readErr == io.EOF {
					return
				}
				log.Printf("Erro ao ler o pacote RTP: %v\n", readErr)
				return
			}

			// Escreve o pacote RTP no buffer
			if err = videoWriter.WriteRTP(packet); err != nil {
				log.Println("Erro ao escrever rtp:", err)
				return
			}

			// Escreve o buffer no arquivo
			_, err = file.Write(buf.Bytes())
			if err != nil {
				log.Println("Erro ao escrever no arquivo", err)
				return
			}
			buf.Reset()
		}
	})

	// Cria um track para enviar o vídeo
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	// Inicia o envio de dados de teste
	go sendVideoData(videoTrack)

	log.Printf("Iniciando a transmissão: %v", stream)
}

func sendVideoData(track *webrtc.TrackLocalStaticRTP) {
	// Configurações para o envio de um vídeo de teste
	ticker := time.NewTicker(time.Second / 30) // Envia um frame a cada 30 fps
	defer ticker.Stop()
	for range ticker.C {
		// Cria um frame de vídeo de teste (pode ser substituído pela captura da câmera)
		frame := createTestFrame()

		if err := track.WriteRTP(frame); err != nil {
			log.Println("Erro ao escrever o frame de teste:", err)
			return
		}
	}
}
func createTestFrame() *rtp.Packet {
	//func createTestFrame() *webrtc.RTPPacket {
	// Cria um frame de teste simples com pixels em formato YUV420
	payload := []byte{
		0x10, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	return &rtp.Packet{
		Header: rtp.Header{
			PayloadType: 100,
		},
		Payload: payload,
	}
	// }
}
