package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	commonv1 "github.com/YangYuS8/phanes/gen/go/phanes/common/v1"
	runtimev1 "github.com/YangYuS8/phanes/gen/go/phanes/runtime/v1"
	"github.com/YangYuS8/phanes/internal/config"
	"github.com/YangYuS8/phanes/internal/contracts"
	"github.com/YangYuS8/phanes/internal/storage"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	cfg       config.Config
	cacheID   string
	startedAt time.Time
	problems  []*commonv1.Problem

	httpServer *http.Server
	listener   net.Listener
	done       chan struct{}
	onReady    func(*Server)

	mu    sync.RWMutex
	state runtimev1.RuntimeState
}

type Options struct {
	OnReady func(*Server)
}

func NewServer(cfg config.Config, opts Options) (*Server, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	verification, err := contracts.VerifyCacheRoot(cfg.Cache.Dir, contracts.CacheVerifyOptions{})
	if err != nil {
		return nil, fmt.Errorf("verify cache before runtime start: %w", err)
	}

	saveDB, err := storage.OpenReadWrite(cfg.Save.DBPath)
	if err != nil {
		return nil, err
	}
	defer saveDB.Close()
	if err := storage.MigrateSave(context.Background(), saveDB); err != nil {
		return nil, err
	}

	server := &Server{
		cfg:       cfg,
		cacheID:   verification.Manifest.GetCacheId(),
		startedAt: time.Now().UTC(),
		state:     runtimev1.RuntimeState_RUNTIME_STATE_STARTING,
		done:      make(chan struct{}),
		onReady:   opts.OnReady,
	}
	server.httpServer = &http.Server{Handler: server.routes()}
	return server, nil
}

func (s *Server) Run(ctx context.Context) error {
	bindHost, err := contracts.NormalizeBindHost(s.cfg.Runtime.BindHost)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", bindHost, s.cfg.Runtime.HTTPPort))
	if err != nil {
		return err
	}
	s.listener = listener
	addr := listener.Addr().(*net.TCPAddr)
	s.cfg.Runtime.HTTPPort = uint32(addr.Port)
	s.setState(runtimev1.RuntimeState_RUNTIME_STATE_READY)
	if s.onReady != nil {
		s.onReady(s)
	}

	go func() {
		<-ctx.Done()
		_ = s.Shutdown(context.Background())
	}()

	err = s.httpServer.Serve(listener)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.setState(runtimev1.RuntimeState_RUNTIME_STATE_FAILED)
		close(s.done)
		return err
	}
	close(s.done)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.setState(runtimev1.RuntimeState_RUNTIME_STATE_STOPPING)
	err := s.httpServer.Shutdown(ctx)
	s.setState(runtimev1.RuntimeState_RUNTIME_STATE_STOPPED)
	return err
}

func (s *Server) Done() <-chan struct{} { return s.done }

func (s *Server) URL() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fmt.Sprintf("http://%s:%d", s.cfg.Runtime.BindHost, s.cfg.Runtime.HTTPPort)
}

func (s *Server) Status() *runtimev1.RuntimeStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &runtimev1.RuntimeStatus{
		State:     s.state,
		BindHost:  s.cfg.Runtime.BindHost,
		HttpPort:  s.cfg.Runtime.HTTPPort,
		CacheId:   s.cacheID,
		Problems:  append([]*commonv1.Problem(nil), s.problems...),
		StartedAt: timestamppb.New(s.startedAt),
	}
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", s.handleHealthz)
	mux.HandleFunc("GET /v1/status", s.handleStatus)
	mux.HandleFunc("GET /v1/diagnostics", s.handleDiagnostics)
	mux.HandleFunc("POST /v1/runtime/shutdown", s.handleShutdown)
	mux.HandleFunc("POST /v1/session/start", s.handleSessionNotImplemented)
	mux.HandleFunc("POST /v1/session/stop", s.handleSessionNotImplemented)
	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeProtoJSON(w, &runtimev1.HealthzResponse{Ok: true, State: s.Status().GetState()})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	writeProtoJSON(w, s.Status())
}

func (s *Server) handleDiagnostics(w http.ResponseWriter, r *http.Request) {
	writeProtoJSON(w, &runtimev1.DiagnosticsResponse{Problems: s.Status().GetProblems()})
}

func (s *Server) handleShutdown(w http.ResponseWriter, r *http.Request) {
	writeProtoJSON(w, &runtimev1.StopRuntimeResult{Accepted: true})
	go func() {
		_ = s.Shutdown(context.Background())
	}()
}

func (s *Server) handleSessionNotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "session endpoints are not implemented yet"})
}

func (s *Server) setState(state runtimev1.RuntimeState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
}

func writeProtoJSON(w http.ResponseWriter, message proto.Message) {
	w.Header().Set("Content-Type", "application/json")
	data, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_, _ = w.Write(data)
}
