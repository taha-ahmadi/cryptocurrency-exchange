import { useState, useEffect } from 'react';
import Web3 from 'web3';
import {ethers} from 'ethers';

const OrderBook = ({ orderData }) => {
  const { TotalAsksVolume, TotalBidsVolume, Asks, Bids } = orderData;

 
  const formatDate = (timestamp) => {
    const dateTime = new Date(timestamp / 1000000); // divide by 1 million for microseconds
    // Format the date and time as desired using the Date object methods
    const formattedDate = `${dateTime.toLocaleDateString()} ${dateTime.toLocaleTimeString()}`;
    return formattedDate;
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 my-5">
    <h2 className="text-3xl font-bold mb-4 text-center">Order Book</h2>
    <div className="grid grid-cols-1 md:grid-cols-2 gap-8 my-5">
      <div className="bg-red-200 rounded-lg p-6 flex items-center justify-center">
        <div>
          <p className="text-2xl font-semibold mb-2 text-red-800">
            Total Asks Volume
          </p>
          <p className="text-xl font-semibold text-center">{TotalAsksVolume}</p>
        </div>
      </div>
      <div className="bg-green-200 rounded-lg p-6 flex items-center justify-center">
        <div>
          <p className="text-2xl font-semibold mb-2 text-green-800">
            Total Bids Volume
          </p>
          <p className="text-xl font-semibold text-center">{TotalBidsVolume}</p>
        </div>
      </div>
    </div>


      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <h3 className="text-xl font-bold mb-4">Asks</h3>
          <table className="w-full border-collapse">
            <thead>
              <tr className="bg-gray-100">
                <th className="text-left py-2 px-3">UserID</th>
                <th className="text-left py-2 px-3">ID</th>
                <th className="text-left py-2 px-3">Amount</th>
                <th className="text-left py-2 px-3">Type</th>
                <th className="text-left py-2 px-3">Price</th>
                <th className="text-left py-2 px-3">Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {Asks.map((ask) => (
                <tr key={ask.ID} className="border-b border-gray-300">
                  <td className="py-2 px-3">{ask.UserID}</td>
                  <td className="py-2 px-3">{ask.ID}</td>
                  <td className="py-2 px-3">{ask.Amount}</td>
                  <td className="py-2 px-3 text-red-600 font-semibold">
                    Ask
                  </td>
                  <td className="py-2 px-3">{ask.Price}</td>
                  <td className="py-2 px-3">{formatDate(ask.Timestamp)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div>
          <h3 className="text-xl font-bold mb-4">Bids</h3>
          <table className="w-full border-collapse">
            <thead>
              <tr className="bg-gray-100">
                <th className="text-left py-2 px-3">UserID</th>
                <th className="text-left py-2 px-3">ID</th>
                <th className="text-left py-2 px-3">Amount</th>
                <th className="text-left py-2 px-3">Type</th>
                <th className="text-left py-2 px-3">Price</th>
                <th className="text-left py-2 px-3">Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {Bids.map((bid) => (
                <tr key={bid.ID} className="border-b border-gray-300">
                  <td className="py-2 px-3">{bid.UserID}</td>
                  <td className="py-2 px-3">{bid.ID}</td>
                  <td className="py-2 px-3">{bid.Amount}</td>
                  <td className="py-2 px-3 text-green-600 font-semibold">
                    Bid
                  </td>
                  <td className="py-2 px-3">{bid.Price}</td>
                  <td className="py-2 px-3">{formatDate(bid.Timestamp)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

const MyPage = () => {
  const [orderData, setOrderData] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch('http://localhost:3000/books/ETH');
      const data = await response.json();
      setOrderData(data);
    };

    fetchData();
  }, []);

  return (
    <div>
      {orderData ? <OrderBook orderData={orderData} /> : <p>Loading...</p>}
    </div>
  );
};

export default function Home() {
  const [walletBalance, setWalletBalance] = useState(null);

  function handleConnect(balance) {
    setWalletBalance(balance);
  }

  return (
    <div>
      <Header onConnect={handleConnect} />
      <div className="container mx-auto mt-8">
        <div className="grid grid-cols-3 gap-4">
          <div className="bg-gradient-to-r from-green-500 to-blue-500 rounded-lg p-6 shadow-md">
            <h2 className="text-lg font-medium mb-4 text-white">Wallet Balance</h2>
            {walletBalance !== null ? (
              <>
                <p className="text-gray-100 text-sm">ETH: {ethers.utils.formatEther(walletBalance)}</p>
                <p className="text-gray-100 text-sm">BTC: 0.12345678</p>
                <p className="text-gray-100 text-sm">USDT: 1000.00</p>
              </>
            ) : null}
          </div>

          {/* Crypto Currency Ask and Bids */}
          <div className="">
            
          </div>

          {/* Market Order */}
          <div className="bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg p-6 shadow-md">
            <h2 className="text-lg font-medium mb-4 text-white">Market Order</h2>
            <form>
              <label htmlFor="amount" className="block text-gray-100 text-sm font-medium mb-2">
                Amount:
              </label>
              <input
                type="text"
                id="amount"
                name="amount"
                placeholder="Enter amount here"
                className="border border-gray-400 rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:border-purple-500 mb-4"
              />
              <button
                type="submit"
                className="bg-white hover:bg-gray-100 text-gray-800 font-bold py-2 px-4 rounded w-full"
              >
                Buy
              </button>
            </form>
          </div>
        </div>
      </div>
      <MyPage />
    </div>
  );
};


const Header = ({ onConnect }) => {
  return (
    <header className="bg-gradient-to-r from-purple-500 to-pink-500">
      <div className="container mx-auto flex justify-between items-center py-4">
        <h1 className="text-white font-bold text-2xl">Awesome Exchange</h1>
        <ConnectButton onConnect={onConnect} />
      </div>
    </header>
  );
};

const ConnectButton = ({ onConnect }) => {
  const [isMetaMaskInstalled, setIsMetaMaskInstalled] = useState(false);
  const [accounts, setAccounts] = useState([]);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    setIsMetaMaskInstalled(typeof window.ethereum !== "undefined");
  }, []);

  useEffect(() => {
    async function getBalance() {
      if (isConnected && accounts.length > 0) {
        try {
          const provider = new ethers.providers.Web3Provider(window.ethereum);
          const balance = await provider.getBalance(accounts[0]);
          onConnect(balance);
        } catch (error) {
          console.error(error);
        }
      }
    }

    getBalance();
  }, [accounts, isConnected, onConnect]);

  async function connectToMetaMask() {
    try {
      await window.ethereum.request({ method: "eth_requestAccounts" });
      const newAccounts = await window.ethereum.request({ method: "eth_accounts" });
      setAccounts(newAccounts);
      setIsConnected(true);
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <div>
      {isMetaMaskInstalled ? (
        <button
          onClick={connectToMetaMask}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          {isConnected ? "Connected!" : "Connect MetaMask!"}
        </button>
      ) : (
        <p>Please install MetaMask to use this feature.</p>
      )}
    </div>
  );
};


const Body = () => {
  const [walletBalance, setWalletBalance] = useState(null);

  function handleConnect(balance) {
    setWalletBalance(balance);
  }

  return (
    <div className="container mx-auto mt-8">
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-gradient-to-r from-green-500 to-blue-500 rounded-lg p-6 shadow-md">
          <h2 className="text-lg font-medium mb-4 text-white">Wallet Balance</h2>
          {walletBalance !== null ? (
            <>
              <p className="text-gray-100 text-sm">ETH: {ethers.utils.formatEther(walletBalance)}</p>
            </>
        ) : (
          <div></div> 
        )}
        </div>
        {/* Crypto Currency Ask and Bids */}
        <div className="bg-gradient-to-r from-yellow-400 via-red-500 to-pink-500 rounded-lg p-6 shadow-md">
          <h2 className="text-lg font-medium mb-4 text-white">Crypto Currency Ask and Bids</h2>
          <table className="w-full table-fixed">
            <thead>
              <tr>
                <th className="text-left w-1/2 px-2 py-2 text-gray-100 text-sm">Ask Price</th>
                <th className="text-right w-1/2 px-2 py-2 text-gray-100 text-sm">Bid Price</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td className="text-left px-2 py-2 text-white text-sm">$4300</td>
                <td className="text-right px-2 py-2 text-white text-sm">$4298</td>
              </tr>
              <tr>
                <td className="text-left px-2 py-2 text-white text-sm">$4301</td>
                <td className="text-right px-2 py-2 text-white text-sm">$4297</td>
              </tr>
              <tr>
                <td className="text-left px-2 py-2 text-white text-sm">$4302</td>
                <td className="text-right px-2 py-2 text-white text-sm">$4296</td>
              </tr>
            </tbody>
          </table>
        </div>

        {/* Market Order */}
        <div className="bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg p-6 shadow-md">
          <h2 className="text-lg font-medium mb-4 text-white">Market Order</h2>
          <form>
            <label htmlFor="amount" className="block text-gray-100 text-sm font-medium mb-2">
              Amount:
            </label>
            <input
              type="text"
              id="amount"
              name="amount"
              placeholder="Enter amount here"
              className="border border-gray-400 rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:border-purple-500 mb-4"
            />
            <button
              type="submit"
              className="bg-white hover:bg-gray-100 text-gray-800 font-bold py-2 px-4 rounded w-full"
            >
              Buy
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};