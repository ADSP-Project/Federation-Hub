import ShopTile from '../ShopTile';
import { ShopTileContainer } from "./MainPage.styles"

const MainPage = ({ shops }) => (
  <ShopTileContainer>
    {shops && shops.map(shop => <ShopTile key={shop.id} shop={shop} />)}
  </ShopTileContainer>
);

export default MainPage;
