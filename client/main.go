package main

import (
	"context"
	"log"
	"time"

	pb "juglar/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Conectamos al servidor gRPC usando grpc.NewClient
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor: %v", err)
	}
	defer conn.Close()

	client := pb.NewJuglarServiceClient(conn)

	// Preparamos el contexto y el mensaje
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req := &pb.MensajeRequest{
		Contenido: "¡Hola Juglar, este es mi primer mensaje por gRPC!",
		Remitente: "Desarrollador",
	}

	// Enviamos el mensaje al servidor
	res, err := client.EnviarMensaje(ctx, req)
	if err != nil {
		log.Fatalf("Error al enviar mensaje: %v", err)
	}

	log.Printf("Respuesta del Servidor -> Status: %s", res.Status)
}
