import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Header from './components/Header';
import ShopList from './components/ShopList';
import MainPage from './components/MainPage';
import GlobalStyles from './globalStyles';
import React, { useState, useEffect } from 'react';

function App() {
  
  const [shops, setShops] = useState([]);

  useEffect(() => {
    fetch('http://localhost:8000/shops')
      .then(response => response.json())
      .then(data => setShops(data))
      .catch((error) => {
        console.error('Error:', error);
      });
  }, []);

  return (
    <Router>
      <div className="App">
        <GlobalStyles />
        <Header />
        <Layout>
          <Routes>
            <Route path="/partners" element={<ShopList shops={shops} />}/>
            <Route path="/" element={<MainPage shops={shops} />} />
          </Routes>
        </Layout>
      </div>
    </Router>
  );
}

export default App;
