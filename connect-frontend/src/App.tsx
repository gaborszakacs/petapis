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
  const [pets, setPets] = useState<string[] >([]);
  return <>
    <ul>
      {pets.map((pet, index) => (
          <li key={index}>
            {pet}
          </li>
      ))}
    </ul>

    <form onSubmit={async (e) => {
      e.preventDefault();
      setInputValue("");
      const response = await client.getPet({
        petId: inputValue,
      });
      setPets((prev) => [
        ...prev,
         `#${response.pet!.petId} ${response.pet!.name}`
      ]);
    }}>
      <input value={inputValue} onChange={e => setInputValue(e.target.value)} />
      <button type="submit">Get Pet</button>
    </form>

    <form onSubmit={async (e) => {
      e.preventDefault();
      for await (const response of client.getPetStream({})) {
        setPets((prev) => [
          ...prev,
          `#${response.pet!.petId} ${response.pet!.name}`
        ]);
      }
      console.log("Stream ended");
    }}>
      <button type="submit">Get Pet Stream</button>
    </form>
  </>;
}

export default App
