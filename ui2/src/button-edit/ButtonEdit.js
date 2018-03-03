import * as React from "react";
import NameEditor from "./NameEditor";
import {Fragment} from "react";

// Shows a list of editors that change to suit the animation
export default class Button extends React.Component {
    constructor(props) {
        super(props);

        // Derive our update state from initial button state
        let button = props.history.location.state.button;
        this.state = {name: button.name}
    }

    render() {
        return <Fragment>
            <NameEditor name={this.state.name} onNameChanged={name => this.setState({name: name})}/>
        </Fragment>
    }
}
