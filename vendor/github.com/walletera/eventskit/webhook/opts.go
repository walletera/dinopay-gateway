package webhook

import "log/slog"

type Opt func(server *Server)

func WithLogger(logger *slog.Logger) Opt {
    return func(server *Server) {
        server.logger = logger
    }
}
