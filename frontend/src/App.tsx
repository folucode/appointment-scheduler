import { useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";
import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient, type Transport } from "@connectrpc/connect";
import { UserService } from "./gen/user_pb";

const transport: Transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

const client = createClient(UserService, transport);

type User = {
  id: string | undefined;
  name: string | undefined;
  email: string | undefined;
};

function App() {
  const [count, setCount] = useState(0);
  const [user, setUser] = useState<User | null>(null);

  const getUser = async () => {
    const response = await client.getUser({
      id: "12345",
    });

    setUser({
      id: response.user?.id,
      name: response.user?.name,
      email: response.user?.email,
    });
  };

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <button onClick={getUser}>get user</button>
        <p>{user ? `${user.name} (${user.email})` : "No user"}</p>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;
