// import { useState } from 'react';
// import apiRequest from '../utils/api';

// const ContactPage = () => {
//   const [name, setName] = useState('');
//   const [responseMessage, setResponseMessage] = useState<string | null>(null);

//   const handleSubmit = async (event: React.FormEvent) => {
//     event.preventDefault();

//     try {
//       const result = await apiRequest<{ message: string }>('/api/v1/user', 'POST', { name });
//       setResponseMessage(result.message);
//     } catch (error) {
//       setResponseMessage('Failed to submit user');
//     }
//   };

//   return (
//     <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4">
//       <h1 className="text-3xl font-bold mb-6">Contact Us</h1>
//       <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
//         <label className="block mb-4">
//           <span className="text-gray-700">Name:</span>
//           <input
//             type="text"
//             value={name}
//             onChange={(e) => setName(e.target.value)}
//             className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
//             required
//           />
//         </label>
//         <button
//           type="submit"
//           className="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
//         >
//           Submit
//         </button>
//       </form>
//       {responseMessage && (
//         <p className="mt-4 text-lg font-semibold text-green-600">{responseMessage}</p>
//       )}
//     </div>
//   );
// };

// export default ContactPage;
