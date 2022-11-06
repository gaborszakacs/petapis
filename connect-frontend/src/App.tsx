import { useState } from 'react'
import './App.css'

import {
  createConnectTransport,
  createPromiseClient,
} from "@bufbuild/connect-web";

// Import service definition that you want to connect to.
import { PetStoreService } from "@buf/bufbuild_connect-web_gaborszakacs-bitrise_pet/pet/v1/pet_connectweb";

// The transport defines what type of endpoint we're hitting.
// In our example we'll be communicating with a Connect endpoint.
const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
  useBinaryFormat: false,
});

// Here we make the client itself, combining the service
// definition with the transport.
const client = createPromiseClient(PetStoreService, transport);

function App() {
  const [inputValue, setInputValue] = useState("");
  const [messages, setMessages] = useState<string[] >([]);
  return <>
    <ol>
      {messages.map((msg, index) => (
          <li key={index}>
            {msg}
          </li>
      ))}
    </ol>
    <form onSubmit={async (e) => {
      e.preventDefault();
      setInputValue("");
      const response = await client.getPet({
        petId: inputValue,
      });
      setMessages((prev) => [
        ...prev,
         `#${response.pet!.petId} ${response.pet!.name}`
      ]);
    }}>
      <input value={inputValue} onChange={e => setInputValue(e.target.value)} />
      <button type="submit">Get Pet</button>
    </form>
  </>;
}

export default App
