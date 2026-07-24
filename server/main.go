package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	pb "juglar/proto"

	"google.golang.org/grpc"
)

var startTime = time.Now()
var totalRequests int64

type MensajeRecord struct {
	Remitente string    `json:"remitente"`
	Contenido string    `json:"contenido"`
	Timestamp time.Time `json:"timestamp"`
}

type server struct {
	pb.UnimplementedJuglarServiceServer
	mu sync.Mutex
}

func (s *server) EnviarMensaje(ctx context.Context, req *pb.MensajeRequest) (*pb.MensajeResponse, error) {
	atomic.AddInt64(&totalRequests, 1)
	log.Printf("Mensaje recibido de [%s]: %s", req.Remitente, req.Contenido)

	nuevoMensaje := MensajeRecord{
		Remitente: req.Remitente,
		Contenido: req.Contenido,
		Timestamp: time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var historial []MensajeRecord
	if archivo, err := os.ReadFile("historial.json"); err == nil {
		json.Unmarshal(archivo, &historial)
	}
	historial = append(historial, nuevoMensaje)

	datos, _ := json.MarshalIndent(historial, "", "  ")
	os.WriteFile("historial.json", datos, 0644)

	// Ajustado estrictamente al campo "Status" definido en el proto limpio
	return &pb.MensajeResponse{
		Status: "OK",
	}, nil
}

func (s *server) ObtenerHistorial(ctx context.Context, req *pb.HistorialRequest) (*pb.HistorialResponse, error) {
	atomic.AddInt64(&totalRequests, 1)
	s.mu.Lock()
	defer s.mu.Unlock()

	var historial []MensajeRecord
	if archivo, err := os.ReadFile("historial.json"); err == nil {
		json.Unmarshal(archivo, &historial)
	}

	var items []*pb.MensajeItem
	for _, m := range historial {
		items = append(items, &pb.MensajeItem{
			Remitente: m.Remitente,
			Contenido: m.Contenido,
			Timestamp: m.Timestamp.Format(time.RFC3339),
		})
	}

	return &pb.HistorialResponse{Mensajes: items}, nil
}

func (s *server) HealthCheck(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	atomic.AddInt64(&totalRequests, 1)
	uptime := time.Since(startTime).Round(time.Second).String()

	return &pb.HealthResponse{
		Estado: "ONLINE",
		Uptime: uptime,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo iniciar el listener: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterJuglarServiceServer(s, &server{})

	log.Println("Juglar Server corriendo en el puerto 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en el servidor: %v", err)
	}
}
