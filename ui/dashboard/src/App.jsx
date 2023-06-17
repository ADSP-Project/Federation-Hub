import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Header from './components/Header';
import ShopList from './components/ShopList';
import MainPage from './components/MainPage';
import GlobalStyles from './GlobalStyles';

function App() {
  return (
    <Router>
      <div className="App">
        <GlobalStyles />
        <Header />
        <Layout>
          <Routes>
            <Route path="/partners" element={<ShopList />} />
            <Route path="/" element={<MainPage />} />
          </Routes>
        </Layout>
      </div>
    </Router>
  );
}

export default App;
