// pages/_app.js
import { Provider } from 'react-redux';
import { PersistGate } from 'redux-persist/integration/react';
import { store, persistor } from '../redux/store';
import SocketProvider from '@/components/SocketProvider';
import '../styles/globals.css';
import Router from '../components/Router';
import { Toaster } from '../components/ui/sonner';  // Make sure this path is correct

function MyApp({ Component, pageProps }) {
  return (
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
      <Toaster />
        <SocketProvider>
          <Router>
            <Component {...pageProps} />
              {/* Add the Toaster component here */}
          </Router>
        </SocketProvider>
      </PersistGate>
    </Provider>
  );
}

export default MyApp;
