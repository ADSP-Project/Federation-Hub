import { Table, Button } from "./ShopList.styles";

const ShopList = ({ shops }) => (
  <Table>
    <thead>
      <tr>
        <th>ID</th>
        <th>Name</th>
        <th>Webhook</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {shops.map(shop => (
        <tr key={shop.id}>
          <td>{shop.id}</td>
          <td>{shop.name}</td>
          <td>{shop.webhookURL}</td>
          <td>
            <Button>Join</Button>
            <Button>Remove</Button>
          </td>
        </tr>
      ))}
    </tbody>
  </Table>
);


export default ShopList;
