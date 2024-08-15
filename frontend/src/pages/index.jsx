import MainLayout from '@/components/MainLayout';
import ProtectedRoutes from '../components/ProtectedRoutes'; // Make sure this component is correctly implemented

const Home = () => {
  return (
    <ProtectedRoutes>
      <MainLayout>
        <h1>Home Page</h1>
      </MainLayout>
    </ProtectedRoutes>
  );
};

export default Home;





// import { useEffect, useState } from 'react';
// import apiRequest from '../utils/api';
// import User from '@/pages/components/user';
// import Link from 'next/link';

// type Data = {
//   name: string;
// };

// const HomePage = () => {
//   const [data, setData] = useState<Data | null>(null);
//   const [error, setError] = useState<string | null>(null);

//   useEffect(() => {
//     const fetchData = async () => {
//       try {
//         const result = await apiRequest<Data>('/api/v1/user');
//         setData(result);
//       } catch (error) {
//         setError('Failed to fetch data');
//       }
//     };

//     fetchData();
//   }, []);

//   return (
//     <div className="min-h-screen bg-gray-100 p-4">
//       <header className="bg-white shadow-md p-4 mb-6 rounded-md">
//         <h1 className="text-3xl font-bold text-gray-800">Users</h1>
//         <Link className="text-indigo-600 hover:text-indigo-800 font-semibold" href="/contact">
//           Contact us
//         </Link>
//       </header>
//       <main className="bg-white p-6 rounded-md shadow-lg max-w-4xl mx-auto">
//         {error && <p className="text-red-500 text-center mb-4">{error}</p>}
//         <User data={data} error={error} />
//       </main>
//       <nav className="mt-6 flex justify-center space-x-4">
//         {/* Add navigation links or buttons here if needed */}
//       </nav>
//     </div>
//   );
// };

// export default HomePage;
