import React from "react";
import { API_URL } from "./constants";
const CreateRoom = (props) => {
    const create = async (e) => {
        e.preventDefault();

        const resp = await fetch(`${API_URL}/create`);
        const { room_id } = await resp.json();

		props.history.push(`/room/${room_id}`)
    };

    return (
        <div>
            <button onClick={create}>Create Room</button>
        </div>
    );
};

export default CreateRoom;
