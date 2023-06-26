import React, { useState } from "react";
import { BrowserRouter, Route, Switch } from "react-router-dom";

import CreateRoom from "./components/CreateRoom"
import Room from "./components/Room"
import WebSocketClient from "./components/WebSocketClient";

function App() {
    return <div className="App">
		<BrowserRouter>
			<Switch>
				<Route path="/" exact component={CreateRoom}></Route>
				<Route path="/room/:roomID" component={Room}></Route>
				<Route path="/client/:userID" component={WebSocketClient}></Route>
			</Switch>
		</BrowserRouter>
	</div>;
}

export default App;
