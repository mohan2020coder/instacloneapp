import React from 'react';
import { Outlet } from 'react-router-dom';
import LeftSidebar from './LeftSidebar';

const MainLayout: React.FC = () => {
  return (
    <div>
      <LeftSidebar />
      <div>
        <Outlet />
      </div>
    </div>
  );
};

export default MainLayout;
