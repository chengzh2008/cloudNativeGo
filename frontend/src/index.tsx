import * as React from "react";
import * as ReactDom from "react-dom";
import {Hello} from "./components/hello";

ReactDom.render(
    <div className="container">
        <h1>MyEvents</h1>
        <Hello name="World" />
    </div>,
    document.getElementById("myevents-app")
);
