import { useState, useEffect } from 'react';
import { ethers } from 'ethers';
import Header from '../components/Header';
import OrderBook from '../components/OrderBook';
import { fetchOrderBook, placeMarketOrder } from '../services/api';

const OrderBookContainer = () => {
  const [orderData, setOrderData] = useState(null);

  useEffect(() => {
    const getOrderBook = async () => {
      try {
        const data = await fetchOrderBook();
        setOrderData(data);
      } catch (error) {
        console.error('Failed to fetch order book:', error);
      }
    };

    getOrderBook();
    
    // Refresh order book every 5 seconds
    const interval = setInterval(getOrderBook, 5000);
    
    return () => clearInterval(interval);
  }, []);

  return (
    <div>
      {orderData ? <OrderBook orderData={orderData} /> : <p>Loading...</p>}
    </div>
  );
};

export default function Home() {
  const [walletBalance, setWalletBalance] = useState(null);
  const [amount, setAmount] = useState('');

  function handleConnect(balance) {
    setWalletBalance(balance);
  }

  async function handleBuy() {
    if (!amount) return;
    
    try {
      await placeMarketOrder({
        UserID: 1, // Default user for demo
        IsBid: true,
        Amount: parseFloat(amount)
      });
      
      // Reset form
      setAmount('');
      
      // You might want to refresh the order book here
      
    } catch (error) {
      console.error('Failed to place market buy order:', error);
    }
  }

  return (
    <div>
      <Header onConnect={handleConnect} />
      <div className="container mx-auto mt-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-gradient-to-r from-green-500 to-blue-500 rounded-lg p-6 shadow-md">
            <h2 className="text-lg font-medium mb-4 text-white">Wallet Balance</h2>
            {walletBalance !== null ? (
              <>
                <p className="text-gray-100 text-sm">ETH: {ethers.utils.formatEther(walletBalance)}</p>
                <p className="text-gray-100 text-sm">BTC: 0.12345678</p>
                <p className="text-gray-100 text-sm">USDT: 1000.00</p>
              </>
            ) : (
              <p className="text-gray-100 text-sm">Connect your wallet to view balance</p>
            )}
          </div>

          {/* Spacer for middle column in larger screens */}
          <div className="hidden md:block"></div>

          {/* Market Order */}
          <div className="bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg p-6 shadow-md">
            <h2 className="text-lg font-medium mb-4 text-white">Market Order</h2>
            <div>
              <label htmlFor="amount" className="block text-gray-100 text-sm font-medium mb-2">
                Amount:
              </label>
              <input
                type="text"
                id="amount"
                name="amount"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="Enter amount here"
                className="border border-gray-400 rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:border-purple-500 mb-4"
              />
              <button
                onClick={handleBuy}
                className="bg-white hover:bg-gray-100 text-gray-800 font-bold py-2 px-4 rounded w-full"
              >
                Buy
              </button>
            </div>
          </div>
        </div>
      </div>
      <OrderBookContainer />
    </div>
  );
} 