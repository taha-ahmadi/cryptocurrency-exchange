import React from 'react';

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

export default OrderBook; 