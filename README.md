# Cryptocurrency Exchange

> This is an educational codebase that demonstrates how cryptocurrency exchanges and matching engines work. It's designed to help understand the core concepts and mechanics of order matching in trading systems.

## Project Structure

```
.
├── cmd/                   # Entry points for executable applications
│   └── exchange/          # Main exchange application
│       └── main.go        # Entry point for the exchange
├── internal/              # Private application and library code
│   ├── api/               # API server code
│   ├── config/            # Configuration management
│   ├── matchingengine/    # Core matching engine
│   └── marketmaker/       # Market maker implementation
├── pkg/                   # Library code that can be used by external applications
├── frontend/              # Web UI
│   ├── src/
│   │   ├── components/    # Reusable UI components
│   │   ├── pages/         # Next.js pages
│   │   └── services/      # API client services
│   ├── public/            # Static files
│   └── ...
├── scripts/               # Scripts and tools
├── documents/             # Documentation and design files
├── Dockerfile             # Main application Dockerfile
├── Dockerfile.ganache     # Ethereum test network Dockerfile
├── docker-compose.yml     # Docker Compose configuration
└── README.md              # Project documentation
```

## Components

### Matching Engine (internal/matchingengine)

The core trading system that matches buy and sell orders. It uses:

- **Orderbook**: Main data structure holding orders
- **Limit**: Represents a price level with orders at that price
- **Order**: Individual buy or sell order with amount, price, etc.

### API Server (internal/api)

RESTful API for interacting with the exchange, including:

- Routes for getting order book, placing orders, canceling orders
- User management
- Ethereum integration for wallet functionality

### Market Maker (internal/marketmaker)

A bot that provides liquidity to the market by continuously placing and updating orders on both sides.

### Frontend (frontend)

React/Next.js application that provides a UI for:

- Viewing the order book
- Placing market and limit orders
- Connecting to MetaMask
- Viewing wallet balances

## Getting Started

### Prerequisites

- Go 1.16+
- Node.js 16+
- Docker and Docker Compose

### Running Locally

1. Clone the repository
2. Set up environment variables by copying `app.example.env` to `app.env` and configuring
3. Start the backend:
   ```
   go run cmd/exchange/main.go
   ```
4. Start the frontend:
   ```
   cd frontend
   npm install
   npm run dev
   ```

### Running with Docker

```
docker-compose up
```

This will start:

- The exchange backend on port 3000
- Ganache (Ethereum testnet) on port 8545
- The frontend on port 8080

## API Endpoints

| Endpoint                   | Method | Description                    |
| -------------------------- | ------ | ------------------------------ |
| `/books/{market}`          | GET    | Get order book for a market    |
| `/books/{market}/best/ask` | GET    | Get best ask for a market      |
| `/books/{market}/best/bid` | GET    | Get best bid for a market      |
| `/orders/{userID}`         | GET    | Get orders for a user          |
| `/orders`                  | POST   | Place an order                 |
| `/orders/{id}`             | DELETE | Cancel an order                |
| `/trades/{market}`         | GET    | Get recent trades for a market |

## License

This is an educational project. Feel free to use it for learning purposes.

## Contributing

Contributions are welcome! Please follow the standard Git workflow:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request
