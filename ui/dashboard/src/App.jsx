import Layout from './components/Layout';
import Header from './components/Header';
import ShopList from './components/ShopList';
import GlobalStyles from './GlobalStyles.jsx';

function App() {
  return (
    <div className="App">
      <GlobalStyles />
      <Header />
      <Layout>
        <ShopList />
      </Layout>
    </div>
  );
}

export default App;
