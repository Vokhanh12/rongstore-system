package main

// import (
// 	"log"
// 	"net"
// 	"net/http"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus/promhttp"
// 	"go.uber.org/zap"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"

// 	iampb "server/api/iam/v1"
// 	wire "server/internal/iam"
// 	"server/internal/iam/infrastructure/client"
// 	"server/pkg/config"
// 	"server/pkg/logger"
// 	"server/pkg/metrics"
// 	"server/pkg/observability"
// )

// func main() {
// 	if err := logger.Init(true); err != nil {
// 		log.Fatalf("failed to init logger: %v", err)
// 	}
// 	zlog := logger.L
// 	defer zlog.Sync()

// 	cfg := config.Load()
// 	zlog.Info("config loaded", zap.String("keycloak_url", cfg.KeycloakURL))

// 	metrics.Register()
// 	go func() {
// 		http.Handle("/metrics", promhttp.Handler())
// 		zlog.Info("metrics endpoint started", zap.String("addr", ":9090"))
// 		if err := http.ListenAndServe(":9090", nil); err != nil {
// 			zlog.Fatal("metrics server failed", zap.Error(err))
// 		}
// 	}()

// 	maxRetries := 10
// 	interval := 3 * time.Second
// 	zlog.Info("checking Keycloak readiness", zap.String("url", cfg.KeycloakURL))
// 	kcClient, err := client.InitKeycloakClient(cfg, maxRetries, interval)
// 	if err != nil {
// 		zlog.Fatal("Keycloak is not ready", zap.Error(err))
// 	}
// 	zlog.Info("Keycloak client ready", zap.String("base_url", kcClient.GetBaseURL()))

// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		zlog.Fatal("failed to listen", zap.Error(err))
// 	}

// 	deps, err := wire.InitializeIamHandler()
// 	if err != nil {
// 		zlog.Fatal("failed to initialize IAM deps", zap.Error(err))
// 	}

// 	deps.Keycloak = kcClient

// 	s := grpc.NewServer(
// 		grpc.ChainUnaryInterceptor(
// 			observability.GrpcTraceUnaryInterceptor(),
// 			observability.UnaryServerInterceptor("iam_service", deps.Store, nil, false),
// 		),
// 	)

// 	reflection.Register(s)
// 	iampb.RegisterIamServiceServer(s, deps.Handler)

// 	zlog.Info("gRPC server started", zap.String("addr", ":50051"))
// 	if err := s.Serve(lis); err != nil {
// 		zlog.Fatal("failed to serve", zap.Error(err))
// 	}
// }

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Player data
type Player struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color int     `json:"color"`
}

// Message generic
type Message struct {
	Type    string   `json:"type"`
	ID      string   `json:"id,omitempty"`
	X       float64  `json:"x,omitempty"`
	Y       float64  `json:"y,omitempty"`
	Color   int      `json:"color,omitempty"`
	Players []Player `json:"players,omitempty"`
}

var (
	players   = make(map[string]*Player)
	playersMu sync.Mutex

	clients   = make(map[string]*net.UDPAddr) // key = playerID
	clientsMu sync.Mutex
)

func handleUDP(conn *net.UDPConn) {
	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			//	log.Println(time.Now().Format(time.RFC3339), "read error:", err)
			continue
		}

		var msg Message
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			//	log.Println(time.Now().Format(time.RFC3339), "json unmarshal error from", addr, ":", err)
			continue
		}

		if msg.Type == "player_update" {
			// update memory
			playersMu.Lock()
			players[msg.ID] = &Player{
				ID:    msg.ID,
				X:     msg.X,
				Y:     msg.Y,
				Color: 0xFF42A5F5,
			}
			playersMu.Unlock()

			// register client addr
			clientsMu.Lock()
			clients[msg.ID] = addr
			clientsMu.Unlock()

			log.Printf("%s - Received update from player %s: X=%.2f Y=%.2f Addr=%s\n",
				time.Now().Format(time.RFC3339), msg.ID, msg.X, msg.Y, addr)
		} else {
			log.Printf("%s - Unknown message type from %s: %s\n",
				time.Now().Format(time.RFC3339), addr, msg.Type)
		}
	}
}

func broadcastSnapshots(conn *net.UDPConn) {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 updates/s
	defer ticker.Stop()

	for range ticker.C {
		playersMu.Lock()
		playerList := make([]Player, 0, len(players))
		for _, p := range players {
			playerList = append(playerList, *p)
		}
		playersMu.Unlock()

		snapshot := Message{
			Type:    "snapshot",
			Players: playerList,
		}
		data, err := json.Marshal(snapshot)
		if err != nil {
			log.Println(time.Now().Format(time.RFC3339), "json marshal error:", err)
			continue
		}

		clientsMu.Lock()
		for id, addr := range clients {
			_, err := conn.WriteToUDP(data, addr)
			if err != nil {
				log.Printf("%s - write error to player %s (%s): %v\n",
					time.Now().Format(time.RFC3339), id, addr, err)
			} else {
				log.Printf("%s - Sent snapshot to player %s (%s) with %d players\n",
					time.Now().Format(time.RFC3339), id, addr, len(playerList))
			}
		}
		clientsMu.Unlock()
	}
}

func main() {

	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println(time.Now().Format(time.RFC3339), "UDP server started at 0.0.0.0:8080")

	// handle incoming packets
	go handleUDP(conn)

	// broadcast snapshots
	go broadcastSnapshots(conn)

	// keep server alive
	select {}
}
