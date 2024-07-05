package app

import (
	"context"
	"di_container/internal/api/access"
	"di_container/internal/api/auth"
	"di_container/internal/api/note"
	"di_container/internal/client/db"
	"di_container/internal/client/db/pg"
	"di_container/internal/client/db/transaction"
	"di_container/internal/client/rpc"
	"di_container/internal/closer"
	"di_container/internal/config"
	"di_container/internal/config/env"
	"di_container/internal/repository"
	noteRepository "di_container/internal/repository/note"
	"di_container/internal/service"
	noteService "di_container/internal/service/note"
	"log"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.SwaggerConfig
	tokenConfig   *env.TokenConfigData

	dbClient            db.Client
	txManager           db.TxManager
	noteRepository      repository.NoteRepository
	noteOtherRepository repository.OtherNoteRepository

	noteService service.NoteService
	authService service.AuthService

	noteImpl   *note.Implementation
	authImpl   *auth.Implementation
	accessImpl *access.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("Failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("Failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("Failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("Failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

func (s *serviceProvider) TokenConfig() *env.TokenConfigData {
	if s.tokenConfig == nil {
		cfg, err := env.NewTokenConfig()
		if err != nil {
			log.Fatalf("Failed to get token config: %s", err.Error())
		}

		s.tokenConfig = cfg
	}

	return s.tokenConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to create db client: %s", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("Ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) NoteRepository(ctx context.Context) repository.NoteRepository {
	if s.noteRepository == nil {
		s.noteRepository = noteRepository.NewRepository(s.DBClient(ctx))
	}

	return s.noteRepository
}

func (s *serviceProvider) NoteService(ctx context.Context) service.NoteService {
	if s.noteService == nil {
		s.noteService = noteService.NewService(
			s.NoteRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.noteService
}

func (s *serviceProvider) GetNoteImpl(ctx context.Context, client rpc.OtherServiceClient) *note.Implementation {
	if s.noteImpl == nil {
		s.noteImpl = note.NewImplementation(s.NoteService(ctx), client)
	}

	return s.noteImpl
}

func (s *serviceProvider) GetAuthImpl() *auth.Implementation {
	if s.authImpl == nil {
		tokenConfig := s.TokenConfig()
		s.authImpl = auth.NewImplementation(tokenConfig)
	}

	return s.authImpl
}

func (s *serviceProvider) GetAccessImpl() *access.Implementation {
	if s.accessImpl == nil {
		tokenConfig := s.TokenConfig()
		s.accessImpl = access.NewImplementation(tokenConfig)
	}

	return s.accessImpl
}
