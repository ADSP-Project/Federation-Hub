import ShopTile from '../ShopTile';
import { ShopTileContainer } from "./MainPage.styles"

const shops = [
  { id: 1, name: 'Sock Shop', description: 'Step into comfort and style with our Sock Shop! Offering a vibrant range of socks that speak louder than words, our collection is sure to knock your socks off. From cozy woolen socks for chilly nights to sleek, moisture-wicking athletic pairs, we have got your feet covered for every occasion.', img: 'https://cdn.shopify.com/s/files/1/0283/2698/5837/files/imgonline-com-ua-compressed-kD2lYL2oA3dw_2000x.jpg?v=1614312080' },
  { id: 2, name: 'Shop 2', description: 'Welcome to our Boutique Shop, the one-stop destination for fashion connoisseurs! Our carefully curated selection features unique, high-quality pieces from emerging designers. Step into a world of fashion-forward, innovative style and discover the joy of finding the perfect outfit that tells your unique story.', img: 'https://img.theculturetrip.com/450x/smart/wp-content/uploads/2018/07/inside-moeon-boutique--moeon.jpg' },
];

const MainPage = () => (
  <ShopTileContainer>
    {shops.map(shop => <ShopTile key={shop.id} shop={shop} />)}
  </ShopTileContainer>
);

export default MainPage;
