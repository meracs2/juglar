# 🚀 Juglar - Sistema de Mensajería gRPC con Go y API Gateway

Proyecto desarrollado en **Go (Golang)** que implementa una arquitectura orientada a microservicios mediante **gRPC** y **Protocol Buffers**, complementado con un **API Gateway HTTP** y un cliente web interactivo.

---

## 🏗️ Arquitectura del Proyecto

El sistema está dividido en componentes desacoplados:

1. **Servidor gRPC (`/server`)**: 
   * Gestiona el backend principal de mensajería y health checks.
   * Almacena y persiste los mensajes localmente en formato JSON (`historial.json`).
   * Utiliza exclusión mutua (`sync.Mutex`) y contadores atómicos para garantizar concurrencia segura.
2. **API Gateway (`/gateway`)**:
   * Servidor HTTP intermediario en Go que traduce peticiones REST a llamadas gRPC (`gRPC Client`).
   * Sirve los archivos estáticos del cliente web (`/web`).
   * Incluye middleware de CORS y un simulador de carga/estrés.
3. **Cliente gRPC de prueba (`/client`)**:
   * Script independiente en Go para probar la comunicación directa con el servidor gRPC por consola.
4. **Cliente Web (`/web`)**:
   * Interfaz gráfica minimalista para interactuar en tiempo real con el historial, envío de mensajes y monitoreo de estado.

---

## 📂 Estructura de Directorios

```text
juglar/
├── client/        # Cliente gRPC de prueba en Go
├── gateway/       # API Gateway HTTP en Go
├── proto/         # Contratos gRPC (.proto y stubs generados)
├── server/        # Servidor backend gRPC en Go
├── web/           # Frontend estático (HTML/CSS/JS)
├── go.mod         # Módulo principal de Go
└── historial.json # Persistencia local de mensajes
⚙️ Requisitos Previos
Tener instalado Go (versión 1.26 o superior recomendada).

Tener instalado el compilador de Protocol Buffers (protoc) si necesitás regenerar los stubs.

🚀 Puesta en Marcha (Local)
Para poner a funcionar todo el sistema de punta a punta, abriri tres terminales distintas en la raíz del proyecto (C:\Users\mhmc\juglar):

1. Iniciar el Servidor gRPC (Backend)
Bash
cd server
go run main.go
(El servidor quedará escuchando en el puerto 50051)

2. Iniciar el API Gateway
Bash
cd gateway
go run main.go
(El gateway levantará el servidor web y el proxy en el puerto 8080)

3. Probar el Cliente gRPC (Opcional)
Bash
cd client
go run main.go
🌐 Abrir la Interfaz Web
Entrá a tu navegador y abrí la dirección:
👉 http://localhost:8080

🛠️ Tecnologías Utilizadas
Golang (Lenguaje principal)

gRPC & Protocol Buffers (Comunicación de alta eficiencia)

net/http (API Gateway y archivos estáticos)

HTML5 / CSS / JavaScript (Frontend)


---

Una vez que lo guardes en la raíz, podés subirlo a GitHub con estos comandos:

```bash
git add README.md
git commit -m "docs: agrega README profesional al repositorio"
git push origin main