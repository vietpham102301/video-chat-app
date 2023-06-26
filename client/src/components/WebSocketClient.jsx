
import React, { useEffect } from "react";
import { w3cwebsocket as WebSocket } from "websocket";
import eruda from 'eruda';

const WebSocketClient = (props) => {
    eruda.init();
    useEffect(() => {
        const webSocket = new WebSocket(`wss://192.168.1.2:8000/establish?userID=${props.match.params.userID}`);

        webSocket.onopen = () => {
            console.log("WebSocket connection established");
            // Send a notification request to the server
            // webSocket.send(JSON.stringify({ type: "notification", message: "Hello server!" }));
        };

        webSocket.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log("Received message:", message);

            // Check the message type
            if (message.type === "notification") {
                console.log("Notification received:", message.notification);
                // Handle the notification as needed
            }
        };

        return () => {
            webSocket.close();
        };
    }, []);

    const handleNotifyClick = async () => {
        const notification = "This is a notification";
        const recipientId = "2"; // Specify the recipient user ID

        try {
            const response = await fetch("https://192.168.1.2:8000/notify", {
                method: "POST",
                body: JSON.stringify({
                    recipientId: recipientId,
                    notification: notification,
                }),
            });

            if (response.ok) {
                console.log("Notification sent successfully");
            } else {
                console.log("Failed to send notification");
            }
        } catch (error) {
            console.log("Error sending notification:", error);
        }
    };


    return (
        <div>
            <h2>WebSocket Client</h2>
            <button onClick={handleNotifyClick}>Send Notification</button>
        </div>
    );
};

export default WebSocketClient;
