package server

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//  dmoj [-h] [-p SERVER_PORT] [-c CONFIG] [-l LOG_FILE] [--no-watchdog] [-a API_PORT] [-A API_HOST] [-s] [-k] [-T TRUSTED_CERTIFICATES]
//            [-e ONLY_EXECUTORS | -x EXCLUDE_EXECUTORS] [--no-ansi] [--skip-self-test]
//            server_host [judge_name] [judge_key]

// https://github.com/DMOJ/online-judge/blob/b1f8c9a09b251695f8b7496a7a2040bfda3ddb65/judge/bridge/base_handler.py#L52
// https://github.com/DMOJ/online-judge/blob/b1f8c9a09b251695f8b7496a7a2040bfda3ddb65/judge/bridge/judge_handler.py

type Judge struct {
	conn      net.Conn
	isRunning bool
	languages []string
	problems  []uuid.UUID
}

type Server struct {
	server net.Listener
	judges map[string]*Judge
}

var Handlers = map[string]func(*Judge, *gorm.DB, string) error{
	"handshake": on_handshake,
}

func NewServer() (*Server, error) {
	s := &Server{}
	server, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
	s.server = server
	if err != nil {
		return nil, err
	}

	s.judges = make(map[string]*Judge)

	fmt.Printf("Listening on %s\n", s.server.Addr().String())
	return s, nil
}

func (server *Server) Handle(id string, conn net.Conn, db *gorm.DB) {
	for {
		size_packet, err := readSize(conn)
		if err != nil {
			conn.Close()
			delete(server.judges, id)
			break
		}
		fmt.Printf("Packet with size %d recieved\n", size_packet)

		packet, err := readPacket(conn, size_packet)
		if err != nil {
			conn.Close()
			delete(server.judges, id)
			break
		}

		data, err := deflatePacket(packet)
		if err != nil {
			conn.Close()
			delete(server.judges, id)
			break
		}

		fmt.Printf("Packet: %s\n", data)

		name, err := getPacketName(data)
		if err != nil {
			conn.Close()
			delete(server.judges, id)
			break
		}

		if handler, ok := Handlers[name]; ok {
			err = handler(server.judges[id], db, data)
			if err != nil {
				conn.Close()
				delete(server.judges, id)
				break
			}
		}
	}
}

func (s *Server) Listen(db *gorm.DB) {
	for {
		conn, err := s.server.Accept()

		if err != nil {
			fmt.Println("Failed to Accept: ", err)
			continue
		}

		id := uuid.NewString()
		s.judges[id] = &Judge{conn: conn}

		go s.Handle(id, conn, db)
	}
}

func (s *Server) Close() error {
	return s.server.Close()
}

// def submit(self, id, problem, language, source):
//
//	data = self.get_related_submission_data(id)
//	self._working = id
//	self._no_response_job = threading.Timer(20, self._kill_if_no_response)
//	self.send({
//			'name': 'submission-request',
//			'submission-id': id,
//			'problem-id': problem,
//			'language': language,
//			'source': source,
//			'time-limit': data.time,
//			'memory-limit': data.memory,
//			'short-circuit': data.short_circuit,
//			'meta': {
//					'pretests-only': data.pretests_only,
//					'in-contest': data.contest_no,
//					'attempt-no': data.attempt_no,
//					'user': data.user_id,
//			},
//	})
func (s *Server) Submit() error {
	var judge *Judge = nil

	// Find a judge that can submit
	for {
		for _, j := range s.judges {
			if !j.isRunning {
				j.isRunning = true
				judge = j
			}
		}

		if judge != nil {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}
