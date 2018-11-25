import * as React from "react";
import * as ReactDom from "react-dom";
import {EventListContainer} from "./components/event_list_container";

ReactDom.render(
    <div className="container">
        <h1>MyEvents</h1>
        <EventListContainer eventListURL="http://localhost:8181/events" />
    </div>,
    document.getElementById("myevents-app")
);
