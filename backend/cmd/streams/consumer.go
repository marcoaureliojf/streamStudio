package main

import (
	// Pacote para trabalhar com buffers de bytes
	"encoding/json" // Pacote para trabalhar com JSON
	"fmt"           // Pacote para formatação de strings
	"io"            // Pacote para operações de input/output
	"log"           // Pacote para logs
	"time"          // Pacote para trabalhar com tempo

	"github.com/marcoaureliojf/streamStudio/backend/internal/config"          // Pacote para configurações
	"github.com/marcoaureliojf/streamStudio/backend/internal/database"        // Pacote para o banco de dados
	"github.com/marcoaureliojf/streamStudio/backend/internal/database/models" // Pacote para os modelos de dados
	"github.com/marcoaureliojf/streamStudio/backend/internal/queue"           // Pacote para a fila de mensagens
	"github.com/pion/rtp"                                                     // Pacote para RTP
	"github.com/pion/webrtc/v3"                                               // Pacote para WebRTC

	// Pacote para trabalhar com mídia WebRTC
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter" // Pacote para escrever pacotes RTP em formato IVF
	"github.com/streadway/amqp"                     // Pacote para RabbitMQ
)

// peerConnection: Variável global para armazenar a conexão WebRTC (inicialmente nula)
var peerConnection *webrtc.PeerConnection

// mediaEngine: Variável global para configurar o motor de mídia do WebRTC
var mediaEngine *webrtc.MediaEngine

func mainConsumer() {
	// Carrega as configurações da aplicação
	cfg := config.LoadConfig()
	// Conecta ao banco de dados
	database.Connect(cfg)

	// Cria uma nova conexão com o RabbitMQ
	rabbitMQ, err := queue.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal("Erro ao iniciar RabbitMQ: ", err) // Se houver erro, encerra o programa
	}
	defer rabbitMQ.Close() // Garante que a conexão com o RabbitMQ seja fechada ao sair da função main

	// Configura o engine de mídia WebRTC, incluindo o codec VP8
	mediaEngine = &webrtc.MediaEngine{}
	mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8, ClockRate: 90000, Channels: 0, SDPFmtpLine: ""},
		PayloadType:        100,
	}, webrtc.RTPCodecTypeVideo)

	// Cria uma nova API WebRTC com o engine de mídia configurado
	api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

	// Inicia o consumo de mensagens da fila, passando a função scheduleHandler para processar as mensagens recebidas
	rabbitMQ.Consume(scheduleHandler(api))
}

// scheduleHandler: Função que retorna um handler para processar as mensagens recebidas da fila de agendamentos
func scheduleHandler(api *webrtc.API) func(d amqp.Delivery) {
	return func(d amqp.Delivery) {
		log.Printf("Recebida mensagem da fila %s\n", d.RoutingKey) // Imprime um log da mensagem recebida

		// Cria uma variável schedule do tipo models.Schedule para desserializar a mensagem JSON
		var schedule models.Schedule
		// Tenta desserializar o corpo da mensagem (d.Body) para um objeto Schedule
		err := json.Unmarshal(d.Body, &schedule)
		if err != nil {
			log.Printf("Erro ao desserializar mensagem: %v", err) // Se houver erro, imprime um log e retorna
			return
		}
		log.Printf("Processando o agendamento %v", schedule) // Imprime um log do agendamento que será processado

		// Busca no banco de dados a transmissão relacionada ao agendamento
		stream, err := findStream(schedule.StreamID)
		if err != nil {
			log.Printf("Erro ao buscar transmissão: %v", err) // Se houver erro, imprime um log e retorna
			return
		}

		// Verifica se já existe uma transmissão em andamento (peerConnection não nula)
		if peerConnection != nil {
			log.Println("Já existe uma transmissão em andamento") // Se já houver, imprime um log e retorna
			return
		}

		// Inicia a transmissão WebRTC
		startStream(stream, api)

		log.Printf("Agendamento %v finalizado\n", schedule) // Imprime um log informando que o processamento do agendamento finalizou
		d.Ack(false)                                        // Envia um reconhecimento (ACK) para o RabbitMQ, informando que a mensagem foi processada
	}
}

// findStream: Função para buscar uma transmissão no banco de dados pelo ID
func findStream(id uint) (*models.Stream, error) {
	// Obtém uma instância do banco de dados
	db := database.GetDB()
	// Cria uma variável stream do tipo models.Stream
	var stream models.Stream
	// Busca a transmissão no banco de dados pelo ID
	result := db.First(&stream, id)
	if result.Error != nil {
		return nil, fmt.Errorf("stream não encontrada") // Se não encontrar, retorna um erro
	}
	return &stream, nil // Se encontrar, retorna um ponteiro para o objeto Stream
}

// startStream: Função para iniciar a transmissão WebRTC
func startStream(stream *models.Stream, api *webrtc.API) {
	// Cria uma nova conexão PeerConnection utilizando as configurações WebRTC
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"}, // Define os servidores STUN para NAT traversal
			},
		},
	})
	if err != nil {
		panic(err) // Se houver erro, interrompe a execução do programa
	}

	// Define o evento OnICEConnectionStateChange que é disparado quando o estado da conexão ICE muda
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Estado da conexão ICE mudou para %s\n", connectionState.String()) // Imprime o novo estado da conexão
		if connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateDisconnected {
			peerConnection = nil // Se a conexão falhar ou for desconectada, libera a variável global para permitir novas transmissões
		}
	})

	// Define o evento OnTrack que é disparado quando uma nova faixa de mídia é recebida
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("Recebido track: %s \n", track.ID()) // Imprime o ID da faixa recebida

		// Cria um arquivo temporário para armazenar o vídeo
		videoWriter, err := ivfwriter.New("output.ivf", nil)
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

		for {
			// Leitura do pacote RTP
			rtp, _, readErr := track.ReadRTP()
			if readErr != nil {
				if readErr == io.EOF {
					return
				}
				log.Printf("Erro ao ler o pacote RTP: %v\n", readErr)
				return
			}

			// Escreve o pacote RTP no buffer
			if err = videoWriter.WriteRTP(rtp); err != nil {
				log.Println("Erro ao escrever rtp:", err)
				return
			}
		}
	})

	// Cria um track para enviar o vídeo
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err) // Se houver erro, interrompe a execução do programa
	}

	// Adiciona o track na conexão WebRTC
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err) // Se houver erro, interrompe a execução do programa
	}
	// Inicia a leitura dos pacotes RTCP para a faixa local (videoTrack), que é uma chamada bloqueante.
	go func() {
		for {
			_, _, err := rtpSender.Read(make([]byte, 1500))
			if err != nil {
				return
			}
		}
	}()

	// Inicia o envio de dados de teste
	go sendVideoData(videoTrack)

	log.Printf("Iniciando a transmissão: %v", stream) // Imprime um log indicando que a transmissão está sendo iniciada
}

func sendVideoData(track *webrtc.TrackLocalStaticRTP) {
	// Configurações para o envio de um vídeo de teste
	ticker := time.NewTicker(time.Second / 30) // Envia um frame a cada 30 fps
	defer ticker.Stop()                        // Garante que o ticker seja parado ao sair da função
	for range ticker.C {
		// Cria um frame de vídeo de teste (pode ser substituído pela captura da câmera)
		frame := createTestFrame()

		if err := track.WriteRTP(frame); err != nil {
			log.Println("Erro ao escrever o frame de teste:", err) // Se houver erro, imprime um log
			return
		}
	}
}
func createTestFrame() *rtp.Packet {
	// Cria um frame de teste simples com pixels em formato YUV420
	payload := []byte{
		0x10, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	return &rtp.Packet{
		Header: rtp.Header{
			PayloadType: 100,
			Version:     2,
		},
		Payload: payload,
	}
}
