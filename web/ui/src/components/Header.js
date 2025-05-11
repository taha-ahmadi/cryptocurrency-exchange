import React from 'react';
import ConnectButton from './ConnectButton';

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

export default Header; 