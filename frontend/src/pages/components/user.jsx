import React from 'react';

const Users = ({ data, error }) => {
  return (
    <div className="p-6 max-w-md mx-auto bg-white rounded-xl shadow-md space-y-4">
      <h1 className="text-2xl font-bold text-center text-gray-800">Users</h1>
      {error && <p className="text-red-500 text-center">{error}</p>}
      {data ? (
        <ul className="space-y-2">
          {data.map(user => (
            <li key={user.id} className="p-2 bg-gray-100 rounded-lg shadow-sm text-gray-700">
              {user.name}
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-center text-gray-500">Loading...</p>
      )}
    </div>
  );
};

export default Users;
