import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import AdminPage from './pages/admin/AdminPage';
import UserPage from './pages/UserPage';

function App() {
  return (
    <Router>
      <div style={{ fontFamily: 'Arial, sans-serif', maxWidth: '1200px', margin: '0 auto', padding: '20px' }}>
        <header style={{ marginBottom: '20px', borderBottom: '1px solid #ccc', paddingBottom: '10px' }}>
          <h1>F1 Prediction App</h1>
          <nav>
            <Link to="/">User View</Link> | <Link to="/admin">Admin View</Link>
          </nav>
        </header>

        <Routes>
          <Route path="/" element={<UserPage />} />
          <Route path="/admin" element={<AdminPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
