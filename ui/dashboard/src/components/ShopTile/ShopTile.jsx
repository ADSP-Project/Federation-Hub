import { Tile, TileHeader, TileBody, TileFooter, JoinButton, Image } from './ShopTile.styles';

const ShopTile = ({ shop }) => (
  <Tile>
    <Image src={shop.img} alt={shop.name} />
    <TileHeader>{shop.name}</TileHeader>
    <TileBody>{shop.description}</TileBody>
    <TileFooter>
      <JoinButton>Join</JoinButton>
    </TileFooter>
  </Tile>
);

export default ShopTile;
