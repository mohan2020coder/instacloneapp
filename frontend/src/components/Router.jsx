// src/components/Router.jsx
import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Home from './Home';
import Profile from './Profile';
import EditProfile from './EditProfile';
import ChatPage from './ChatPage';
import Login from './Login';
import Signup from './Signup';
import MainLayout from './MainLayout';
import ProtectedRoutes from './ProtectedRoutes'; // Ensure this is correctly implemented

const Router = () => {
  return (
    <BrowserRouter>
       <Routes>
        <Route path="/" element={<ProtectedRoutes><MainLayout /></ProtectedRoutes>}>
          <Route index element={<ProtectedRoutes><Home /></ProtectedRoutes>} />
          <Route path="profile/:id" element={<ProtectedRoutes><Profile /></ProtectedRoutes>} />
          <Route path="account/edit" element={<ProtectedRoutes><EditProfile /></ProtectedRoutes>} />
          <Route path="chat" element={<ProtectedRoutes><ChatPage /></ProtectedRoutes>} />
        </Route>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
