import * as React from "react";
import './Virtual.css';
import {Open} from "./server-proxy/frameBuffer";

// Displays the frame buffer as a two dimensional array of pixels
export default class Virtual extends React.Component {
    constructor(props) {
        super(props);
        this.state = {leds: []};
    }

    ws = null;

    // Open web socket to stream frame data
    componentDidMount() {
        this.ws = Open(leds => {
            this.setState({leds: leds});
        })
    }

    // Close web socket
    componentWillUnmount() {
        this.ws.close()
    }

    render() {
        let i = 0;
        return <div>
            <table>
                <tbody>{
                    this.state.leds.map(r =>
                        <tr key={i++}>
                            <VirtualRow rowData={r} rowIndex={i}/>
                        </tr>)}</tbody>
            </table>
        </div>
    }
}

// Render a virtual row
function VirtualRow(props) {
    let i = 0;

    // Row index in first column
    let indexColumn = <td key={i++}>
        {props.rowIndex}
    </td>;

    // Disabled text if row not in use in second column
    if (props.rowData.length === 0)
        return [indexColumn, <td key={i++} className="disabled-text">Disabled</td>];

    // Otherwise a pixel per column for an active row
    return [indexColumn, props.rowData.map(c =>
        <td style={{
            backgroundColor: "#" + c,
            width: "10px",
            height: "10px"
        }} key={i++}/>)]
}
