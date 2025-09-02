package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Player struct {
	ID    string  `json:"id"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color int     `json:"color"`
}

type Message struct {
	Type    string   `json:"type"`
	ID      string   `json:"id,omitempty"`
	X       float64  `json:"x,omitempty"`
	Y       float64  `json:"y,omitempty"`
	Color   int      `json:"color,omitempty"`
	Players []Player `json:"players,omitempty"`
}

func main() {
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	conn, _ := net.DialUDP("udp", nil, serverAddr)
	defer conn.Close()

	player1 := "player1"
	player2 := "player2"

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Println("read error:", err)
				continue
			}
			var snapshot Message
			if err := json.Unmarshal(buf[:n], &snapshot); err != nil {
				log.Println("json unmarshal error:", err)
				continue
			}
			fmt.Println("Snapshot received:", snapshot.Players)
		}
	}()

	go func() {
		x, y := 0.0, 0.0
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			msg := Message{
				Type: "player_update",
				ID:   player1,
				X:    x,
				Y:    y,
			}
			data, _ := json.Marshal(msg)
			conn.Write(data)
			x += 1.0
			y += 1.0
		}
	}()

	go func() {
		x, y := 10.0, 10.0
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			msg := Message{
				Type: "player_update",
				ID:   player2,
				X:    x,
				Y:    y,
			}
			data, _ := json.Marshal(msg)
			conn.Write(data)
			x -= 1.0
			y -= 1.0
		}
	}()

	select {}
}
