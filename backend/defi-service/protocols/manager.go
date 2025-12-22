package protocols

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Manager manages all protocol integrations
type Manager struct {
	ethClient  *ethclient.Client
	baseClient *ethclient.Client
	
	protocols map[string]Protocol
	mu        sync.RWMutex
}

// Protocol interface for DeFi protocols
type Protocol interface {
	GetName() string
	GetAPY(ctx context.Context, asset string, chain string) (float64, error)
	GetUserPositions(ctx context.Context, userAddress string, chain string) ([]Position, error)
	GetHealthFactor(ctx context.Context, userAddress string, chain string) (float64, error)
	GetAssetPrice(ctx context.Context, asset string, chain string) (float64, error)
}

// Position represents a DeFi position
type Position struct {
	Protocol     string  `json:"protocol"`
	Chain        string  `json:"chain"`
	Asset        string  `json:"asset"`
	Type         string  `json:"type"` // lending, borrowing, staking
	Amount       float64 `json:"amount"`
	APY          float64 `json:"apy"`
	Address      string  `json:"address"`
}

// NewManager creates a new protocol manager
func NewManager(ethRPC, baseRPC string) *Manager {
	ethClient, err := ethclient.Dial(ethRPC)
	if err != nil {
		panic("Failed to connect to Ethereum: " + err.Error())
	}

	baseClient, err := ethclient.Dial(baseRPC)
	if err != nil {
		panic("Failed to connect to Base: " + err.Error())
	}

	m := &Manager{
		ethClient:  ethClient,
		baseClient: baseClient,
		protocols:  make(map[string]Protocol),
	}

	// Register protocols
	m.RegisterProtocol(NewAave(ethClient, baseClient))
	m.RegisterProtocol(NewCompound(ethClient, baseClient))
	m.RegisterProtocol(NewEigenLayer(ethClient))

	return m
}

// RegisterProtocol registers a protocol
func (m *Manager) RegisterProtocol(p Protocol) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.protocols[p.GetName()] = p
}

// GetProtocol returns a protocol by name
func (m *Manager) GetProtocol(name string) (Protocol, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.protocols[name]
	return p, ok
}

// GetAllProtocols returns all registered protocols
func (m *Manager) GetAllProtocols() []Protocol {
	m.mu.RLock()
	defer m.mu.RUnlock()
	protocols := make([]Protocol, 0, len(m.protocols))
	for _, p := range m.protocols {
		protocols = append(protocols, p)
	}
	return protocols
}

// GetClient returns the appropriate client for a chain
func (m *Manager) GetClient(chain string) *ethclient.Client {
	if chain == "base" {
		return m.baseClient
	}
	return m.ethClient
}

// StartDataRefresh starts background data refresh
func (m *Manager) StartDataRefresh(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Refresh data for all protocols
			// This can be extended to cache APY data, prices, etc.
		}
	}
}

