import React, { useState, useEffect } from 'react';
import { ethers } from 'ethers';

const ConnectButton = ({ onConnect }) => {
  const [isMetaMaskInstalled, setIsMetaMaskInstalled] = useState(false);
  const [accounts, setAccounts] = useState([]);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    setIsMetaMaskInstalled(typeof window.ethereum !== "undefined");
  }, []);

  async function getBalance() {
    try {
      const provider = new ethers.providers.Web3Provider(window.ethereum);
      const address = accounts[0];
      const balance = await provider.getBalance(address);
      onConnect(balance);
    } catch (error) {
      console.error("Error getting balance:", error);
    }
  }

  async function connectToMetaMask() {
    try {
      const accounts = await window.ethereum.request({
        method: "eth_requestAccounts",
      });
      setAccounts(accounts);
      setIsConnected(true);
      getBalance();
    } catch (error) {
      console.error("Error connecting to MetaMask:", error);
    }
  }

  if (!isMetaMaskInstalled) {
    return (
      <a
        href="https://metamask.io/download/"
        target="_blank"
        rel="noopener noreferrer"
        className="bg-white hover:bg-gray-100 text-gray-800 font-bold py-2 px-4 rounded"
      >
        Install MetaMask
      </a>
    );
  }

  if (isConnected) {
    return (
      <button className="bg-white hover:bg-gray-100 text-gray-800 font-bold py-2 px-4 rounded">
        {accounts[0].substring(0, 6)}...
        {accounts[0].substring(accounts[0].length - 4)}
      </button>
    );
  }

  return (
    <button
      onClick={connectToMetaMask}
      className="bg-white hover:bg-gray-100 text-gray-800 font-bold py-2 px-4 rounded"
    >
      Connect Wallet
    </button>
  );
};

export default ConnectButton; 