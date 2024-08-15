// // redux/store.js
// import { combineReducers, configureStore } from "@reduxjs/toolkit";
// import authSlice from "./authSlice.js";
// import postSlice from './postSlice.js';
// import socketSlice from "./socketSlice.js";
// import chatSlice from "./chatSlice.js";
// import rtnSlice from "./rtnSlice.js";

// import { 
//     persistReducer,
//     FLUSH,
//     REHYDRATE,
//     PAUSE,
//     PERSIST,
//     PURGE,
//     REGISTER,
//     persistStore
// } from 'redux-persist';
// import storage from 'redux-persist/lib/storage';

// const persistConfig = {
//     key: 'root',
//     version: 1,
//     storage,
// };

// const rootReducer = combineReducers({
//     auth: authSlice,
//     post: postSlice,
//     socketio: socketSlice,
//     chat: chatSlice,
//     realTimeNotification: rtnSlice,
// });

// const persistedReducer = persistReducer(persistConfig, rootReducer);

// const store = configureStore({
//     reducer: persistedReducer,
//     middleware: (getDefaultMiddleware) =>
//         getDefaultMiddleware({
//             serializableCheck: {
//                 ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
//             },
//         }),
// });

// const persistor = persistStore(store);

// export { store, persistor };
// redux/store.ts

import { combineReducers, configureStore } from '@reduxjs/toolkit';
import { persistReducer, persistStore } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import authSlice from './authSlice';
import postSlice from './postSlice';
import socketSlice from './socketSlice';
import chatSlice from './chatSlice';
import rtnSlice from './rtnSlice';

const persistConfig = {
  key: 'root',
  version: 1,
  storage,
};

const rootReducer = combineReducers({
  auth: authSlice,
  post: postSlice,
  socketio: socketSlice,
  chat: chatSlice,
  realTimeNotification: rtnSlice,
});

const persistedReducer = persistReducer(persistConfig, rootReducer);

const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [
          'persist/FLUSH',
          'persist/REHYDRATE',
          'persist/PAUSE',
          'persist/PERSIST',
          'persist/PURGE',
          'persist/REGISTER',
        ],
      },
    }),
});

const persistor = persistStore(store);

export { store, persistor };
