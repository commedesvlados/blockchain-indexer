package handler

import (
	"context"
	"github.com/commedesvlados/blockchain-indexer/pkg/repository"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Handler struct {
	client     *ethclient.Client
	repository *repository.Repositiry
}

func NewHandler(repository *repository.Repositiry, client *ethclient.Client) *Handler {
	return &Handler{
		repository: repository,
		client:     client,
	}
}

func (h *Handler) StartListen(ctx context.Context) {
	go h.BlockListener(ctx)
	go h.ERC20Listener(ctx)
}
