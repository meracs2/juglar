package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	pb "juglar/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MensajePayload struct {
	Remitente string `json:"remitente"`
	Contenido string `json:"contenido"`
}

type StressResponse struct {
	Peticiones int    `json:"peticiones"`
	DuracionMs int64  `json:"duracion_ms"`
	Estado     string `json:"estado"`
}

// Middleware simple para habilitar CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	// Conexión gRPC persistente y reutilizable
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error al conectar con el servidor gRPC: %v", err)
	}
	defer conn.Close()

	// Cliente gRPC global para la mensajería
	grpcClient := pb.NewJuglarServiceClient(conn)

	// Servidor de archivos estáticos
	fs := http.FileServer(http.Dir("../web"))
	http.Handle("/", fs)

	// Endpoint de Health Check
	http.HandleFunc("/api/health", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		res, err := grpcClient.HealthCheck(ctx, &pb.HealthRequest{})
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"estado":"OFFLINE","uptime":"N/A"}`))
			return
		}

		json.NewEncoder(w).Encode(res)
	}))

	// Endpoint de Simulador de Carga / Estrés gRPC
	http.HandleFunc("/api/stress", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		const totalCalls = 50
		var wg sync.WaitGroup
		var mu sync.Mutex
		successCount := 0
		inicio := time.Now()

		for i := 0; i < totalCalls; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()

				_, err := grpcClient.HealthCheck(ctx, &pb.HealthRequest{})
				if err == nil {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()
		duracion := time.Since(inicio).Milliseconds()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StressResponse{
			Peticiones: successCount,
			DuracionMs: duracion,
			Estado:     "EXITOSO",
		})
	}))

	// Endpoint historial que consulta directo al gRPC
	http.HandleFunc("/api/historial", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		res, err := grpcClient.ObtenerHistorial(ctx, &pb.HistorialRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Mensajes)
	}))

	// Endpoint enviar mensaje
	http.HandleFunc("/api/enviar", enableCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		var payload MensajePayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		req := &pb.MensajeRequest{
			Remitente: payload.Remitente,
			Contenido: payload.Contenido,
		}

		res, err := grpcClient.EnviarMensaje(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}))

	log.Println("Gateway corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
