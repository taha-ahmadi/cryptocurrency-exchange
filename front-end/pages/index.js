import { useState, useEffect } from 'react';
import Web3 from 'web3';
import {ethers} from 'ethers';


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