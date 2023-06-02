import { Table, Button } from "./ShopList.styles";

const shops = [
    { id: 1, name: 'Shop 1', description: 'Description 1' },
    { id: 2, name: 'Shop 2', description: 'Description 2' },
  ];
  
  const ShopList = () => (
    <Table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Description</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {shops.map(shop => (
          <tr key={shop.id}>
            <td>{shop.id}</td>
            <td>{shop.name}</td>
            <td>{shop.description}</td>
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
  